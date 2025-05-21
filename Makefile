deploy:
	gcloud run deploy lysk-battle-record \
      --source . \
      --region asia-southeast1 \
      --allow-unauthenticated \
#      --set-env-vars GOOGLE_APPLICATION_CREDENTIALS=/Users/linfanshi/Documents/lysk-battle-record/credentials.json