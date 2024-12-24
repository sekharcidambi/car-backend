docker build -t car-backend .;
docker tag car-backend:latest gcr.io/carpool-app-backend/car-backend-latest;
docker push gcr.io/carpool-app-backend/car-backend-latest;
gcloud run deploy car-backend-latest \
    --image gcr.io/carpool-app-backend/car-backend-latest \
    --add-cloudsql-instances carpool-app-backend:us-west1:carpool-db \
    --project carpool-app-backend \
    --region us-west1 \
    --service-account="car-backend-sa@carpool-app-backend.iam.gserviceaccount.com" \
    --set-env-vars "INSTANCE_CONNECTION_NAME=carpool-app-backend:us-west1:carpool-db,DB_USER=postgres,DB_NAME=carpool" \
    --set-secrets "DB_PASSWORD=db-password:latest" \
    --min-instances 1 \
    --timeout 300s \
    --cpu 1 \
    --memory 512Mi
