aws --endpoint-url http://seaweedfs:8333 \
  s3api create-bucket \
  --bucket private

aws --endpoint-url http://seaweedfs:8333 \
  s3api create-bucket \
  --bucket public

aws --endpoint-url http://seaweedfs:8333 \
  s3api put-bucket-policy \
  --bucket public \
  --policy file:///scripts/public-bucket-policy.json
