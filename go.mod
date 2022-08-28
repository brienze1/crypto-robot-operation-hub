module github.com/brienze1/crypto-robot-operation-hub

go 1.19

require (
	github.com/aws/aws-lambda-go v1.34.1
	github.com/aws/aws-sdk-go v1.44.85
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.4.0
	github.com/stretchr/testify v1.7.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//intentional vulnerable dependencies for testing github action dep checker
require (
		cloud.google.com/go v0.40.0 // indirect
    	code.gitea.io/gitea v1.9.0-dev
    	gitea.com/lunny/levelqueue v0.1.0 // indirect
    	gitea.com/macaron/inject v0.0.0-20190805023432-d4c86e31027a
    	gitea.com/macaron/macaron v1.3.2
    	github.com/BurntSushi/xgb v0.0.0-20160522181843-27f122750802 // indirect
    	github.com/OneOfOne/xxhash v1.2.5 // indirect
    	github.com/PuerkitoBio/goquery v1.5.0 // indirect
    	github.com/RoaringBitmap/roaring v0.4.7 // indirect
    	github.com/Unknwon/com v0.0.0-20190321035513-0fed4efef755 // indirect
    	github.com/apache/thrift v0.12.0 // indirect
    	github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5 // indirect
    	github.com/bgentry/speakeasy v0.1.0 // indirect
    	github.com/blevesearch/bleve v0.0.0-20190214220507-05d86ea8f6e3 // indirect
    	github.com/blevesearch/blevex v0.0.0-20180227211930-4b158bb555a3 // indirect
    	github.com/blevesearch/go-porterstemmer v0.0.0-20141230013033-23a2c8e5cf1f // indirect
    	github.com/blevesearch/segment v0.0.0-20160105220820-db70c57796cc // indirect
    	github.com/boombuler/barcode v0.0.0-20161226211916-fe0f26ff6d26 // indirect
    	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
    	github.com/chaseadamsio/goorgeous v0.0.0-20170901132237-098da33fde5f // indirect
    	github.com/corbym/gocrest v1.0.3 // indirect
    	github.com/coreos/bbolt v1.3.3 // indirect
    	github.com/coreos/etcd v3.3.15+incompatible // indirect
    	github.com/coreos/go-oidc v2.1.0+incompatible // indirect
)