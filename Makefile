DEPENDENCIES = "github.com/go-sql-driver/mysql" "github.com/gorilla/context" "github.com/gorilla/sessions" "github.com/lestrrat/go-ngram" "github.com/microcosm-cc/bluemonday" "github.com/mrjones/oauth" "github.com/russross/blackfriday" "golang.org/x/oauth2" "golang.org/x/oauth2/google" "github.com/jinzhu/gorm" "github.com/gholt/blackfridaytext"

all: hsproom hsp3dishjs ace

hsproom: go_packages
	gom build hsproom.go

go_packages:
	gom install $(DEPENDENCIES)

hsp3dishjs: submodules
	cd openhsp/hsp3dish/; \
	mv extlib/src/Lua extlib/src/lua; \
	emmake make -f makefile.emscripten; \
	sed -i -e 's/\(var\ webGLContextAttributes[^}]*\)\(};\)/\1,preserveDrawingBuffer:true\2/' emscripten/hsp3dish.js emscripten/hsp3dish-gp.js; \
	cd ../../

ace: submodules
	cd ace/; \
	npm install; \
	node Makefile.dryice.js -m -nc; \
	cd ../

submodules: .git
	git submodule update --init

clean:
	rm -f hsproom
