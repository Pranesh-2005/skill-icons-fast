package handler

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

/* =========================================================
   Embedded icons JSON
   ========================================================= */

//go:embed icons.json
var iconsJSON []byte

var icons = make(map[string]string)
var iconNameList []string
var themedIcons []string

/* =========================================================
   Short name aliases
   ========================================================= */

var shortNames = map[string]string{
	"js":                "javascript",
	"ts":                "typescript",
	"py":                "python",
	"tailwind":          "tailwindcss",
	"vue":               "vuejs",
	"nuxt":              "nuxtjs",
	"go":                "golang",
	"cf":                "cloudflare",
	"wasm":              "webassembly",
	"postgres":          "postgresql",
	"k8s":               "kubernetes",
	"next":              "nextjs",
	"mongo":             "mongodb",
	"md":                "markdown",
	"ps":                "photoshop",
	"ai":                "illustrator",
	"pr":                "premiere",
	"ae":                "aftereffects",
	"scss":              "sass",
	"sc":                "scala",
	"net":               "dotnet",
	"gatsbyjs":          "gatsby",
	"gql":               "graphql",
	"vlang":             "v",
	"amazonwebservices": "aws",
	"bots":              "discordbots",
	"express":           "expressjs",
	"googlecloud":       "gcp",
	"mui":               "materialui",
	"windi":             "windicss",
	"unreal":            "unrealengine",
	"nest":              "nestjs",
	"ktorio":            "ktor",
	"pwsh":              "powershell",
	"au":                "audition",
	"rollup":            "rollupjs",
	"rxjs":              "reactivex",
	"rxjava":            "reactivex",
	"ghactions":         "githubactions",
	"sklearn":           "scikitlearn",
	"ml5":               "ml5js",
	"vb":                "visualbasic",
	"an":                "animate",
	"ca":                "capture",
	"cc":                "creativecloud",
	"ch":                "characteranimator",
	"me":                "mediaencoder",
	"pl":                "prelude",
	"ru":                "premiererush",
	"fs":                "fuse",
	"id":                "indesign",
	"ic":                "incopy",
	"sp":                "adobespark",
	"dw":                "dreamweaver",
	"dn":                "dimension",
	"ar":                "aero",
	"psc":               "photoshopclassic",
	"psx":               "photoshopexpress",
	"lr":                "lightroom",
	"lrc":               "lightroomclassic",
	"fr":                "fresco",
	"pf":                "portfolio",
	"st":                "stock",
	"be":                "behance",
	"br":                "bridge",
	"million":           "millionjs",
	"asm":               "assembly",
	"pop":               "popos",
	"nix":               "nixos",
	"hc":                "holyc",
	"yml":               "yaml",
	"twitter":           "x",
	"arc":               "arcbrowser",
	"hf":                "huggingface",
	"sqla":              "sqlalchemy",
	"notepad++":         "notepadpp",
	"jq":                "jqlang",
}

/* =========================================================
   App setup
   ========================================================= */

var app *gin.Engine

const (
	ICONS_PER_LINE = 15
	ONE_ICON       = 48
	SCALE          = float64(ONE_ICON) / float64(300-44)
)

/* =========================================================
   SVG generation
   ========================================================= */

func generateSvg(iconNames []string, perLine int, hasTitles bool, align string) string {
	iconSvgList := make([]string, 0, len(iconNames))

	for _, name := range iconNames {
		svg, ok := icons[name]
		if !ok {
			continue // ✅ prevent panic
		}
		iconSvgList = append(iconSvgList, svg)
	}

	length := int(math.Min(float64(perLine*300), float64(len(iconSvgList)*300))) - 44
	height := int(math.Ceil(float64(len(iconSvgList))/float64(perLine)))*300 - 44
	scaledHeight := int(float64(height) * SCALE)
	scaledWidth := int(float64(length) * SCALE)

	var svg string
	switch align {
	case "center":
		svg = fmt.Sprintf(`<svg width="100%%" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
			scaledHeight, length, height)
	case "right":
		svg = fmt.Sprintf(`<svg width="calc(200%% - %dpx)" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
			scaledWidth, scaledHeight, length, height)
	default:
		svg = fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
			scaledWidth, scaledHeight, length, height)
	}

	for i, icon := range iconSvgList {
		title := ""
		if hasTitles {
			title = fmt.Sprintf("<title>%s</title>", iconNames[i])
		}
		x := (i % perLine) * 300
		y := (i / perLine) * 300
		svg += fmt.Sprintf(`<g transform="translate(%d,%d)">%s%s</g>`, x, y, title, icon)
	}

	return svg + "</svg>"
}

/* =========================================================
   Helpers
   ========================================================= */

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func parseShortNames(names []string, theme string) []string {
	result := make([]string, 0, len(names))
	for _, name := range names {
		if contains(iconNameList, name) {
			if contains(themedIcons, name) {
				result = append(result, name+"-"+theme)
			} else {
				result = append(result, name)
			}
			continue
		}
		if val, ok := shortNames[name]; ok {
			if contains(themedIcons, val) {
				result = append(result, val+"-"+theme)
			} else {
				result = append(result, val)
			}
		}
	}
	return result
}

/* =========================================================
   Routes
   ========================================================= */

func iconRoute(r *gin.RouterGroup) {
	r.GET("/icons", func(ctx *gin.Context) {
		iconParam := ctx.Query("i")
		if iconParam == "" {
			ctx.String(http.StatusBadRequest, "You didn't specify any icons!")
			return
		}

		theme := ctx.DefaultQuery("theme", "auto")
		perLine, _ := strconv.Atoi(ctx.DefaultQuery("perline", "15"))
		align := ctx.DefaultQuery("align", "left")
		hasTitles := ctx.Query("titles") != ""

		var iconShortNames []string
		if iconParam == "all" {
			iconShortNames = iconNameList
		} else {
			iconShortNames = strings.Split(iconParam, ",")
		}

		iconNames := parseShortNames(iconShortNames, theme)
		svg := generateSvg(iconNames, perLine, hasTitles, align)

		ctx.Header("Content-Type", "image/svg+xml")
		ctx.Header("Cache-Control", "public, max-age=31556952, s-maxage=31536000")
		ctx.String(http.StatusOK, svg)
	})
}

/* =========================================================
   Init
   ========================================================= */

func init() {
	if err := json.Unmarshal(iconsJSON, &icons); err != nil {
		panic(err)
	}

	if len(icons) == 0 {
		panic("icons.json loaded but empty")
	}

	for key := range icons {
		base := strings.Split(key, "-")[0]
		iconNameList = append(iconNameList, base)
		if strings.Contains(key, "-light") ||
			strings.Contains(key, "-dark") ||
			strings.Contains(key, "-auto") {
			themedIcons = append(themedIcons, base)
		}
	}

	app = gin.New()
	app.Use(gin.Recovery()) // ✅ critical for Vercel
	r := app.Group("/api")
	iconRoute(r)
}

/* =========================================================
   Vercel entrypoint
   ========================================================= */

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
