build:
	docker build -t lysk-server .

run:
    docker run -p 8080:8080 lysk-server

deploy:
	gcloud run deploy lysk-battle-record \
      --source . \
      --region asia-southeast1 \
      --allow-unauthenticated \
      --service-account sheets-accessor@double-voice-460107-e4.iam.gserviceaccount.com
