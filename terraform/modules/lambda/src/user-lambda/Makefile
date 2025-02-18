TABLES := devices,users

gen-dao:
	gentool -dsn "root:password@tcp(localhost:3306)/ih_authenticator?charset=utf8mb4&parseTime=True&loc=Local" -fieldNullable -onlyModel --outPath="./infrastructure/mysql/gorm/model" \
	--tables="$(TABLES)"
	find ./infrastructure/mysql/gorm/model -type f -name '*.gen.go' | xargs sed -i '' -e 's/int32/int64/g'
