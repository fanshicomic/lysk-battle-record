package usecases

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"lysk-battle-record/internal/models"
)

func (s *LyskServer) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logrus.Info("[Auth] No authorization header provided, proceeding without authentication")
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logrus.Warnf("[Auth] Invalid authorization header format: %s", authHeader)
			c.Next()
			return
		}

		tokenString := parts[1]
		userID, err := s.auth.ValidateJWT(tokenString)
		if err != nil {
			logrus.Warnf("[Auth] JWT validation failed: %v", err)
			c.Next()
			return
		}

		logrus.Infof("[Auth] Successfully authenticated user: %s", userID)
		c.Set("userID", userID)
		c.Next()
	}
}

func (s *LyskServer) Login(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := s.auth.Login(req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, err := s.auth.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusNonAuthoritativeInfo, gin.H{"error": err.Error()})
		return
	}

	if err := s.createUserIfNotExist(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *LyskServer) GetRanking(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		userId = ""
	}
	ranking := s.orbitRecordStore.GetRanking(userId.(string))

	c.JSON(http.StatusOK, ranking)
}

func (s *LyskServer) createUserIfNotExist(userId string) error {
	var user models.User
	user.ID = userId

	_, ok := s.userStore.Get(userId)
	if ok {
		return nil
	}

	createdUser, err := s.userSheetClient.ProcessUser(user)
	if err != nil {
		return fmt.Errorf("创建用户失败: %v", err)
	}

	s.userStore.Insert(*createdUser)
	return nil
}

func (s *LyskServer) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或无效的用户"})
		return
	}
	user.ID = userId.(string)

	_, ok := s.userStore.Get(userId.(string))
	if ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户已存在"})
		return
	}

	if err := user.ValidateNickname(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := s.userSheetClient.ProcessUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.userStore.Insert(*createdUser)
	c.JSON(http.StatusOK, createdUser)
}

func (s *LyskServer) GetUser(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或无效的用户"})
		return
	}
	user, ok := s.userStore.Get(userId.(string))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (s *LyskServer) UpdateUser(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或无效的用户"})
		return
	}
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user.ID = userId.(string)

	currentUser, ok := s.userStore.Get(userId.(string))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	if err := user.ValidateNickname(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.RowNumber = currentUser.RowNumber

	if err := s.userSheetClient.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := s.userStore.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
