package compound

import (
	"errors"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

var (
	globalFontPath string
	globalFontType *truetype.Font
)

type Compound struct {
	width      float64
	height     float64
	dpi        float64
	fontPath   string
	lineSpace  float64
	obj        *image.RGBA
	fontType   *truetype.Font
	once       sync.Once
	fontDrawer *font.Drawer

	Point *fixed.Point26_6
}

func NewCompound() *Compound {
	this := new(Compound)
	return this
}

func SetCompound(fontPath string) error {
	globalFontPath = fontPath
	fontBytes, err := ioutil.ReadFile(globalFontPath)
	if err != nil {
		return err
	}
	globalFontType, err = truetype.Parse(fontBytes)
	if err != nil {
		return err
	}
	return nil
}

//1200,800,"./src/simsun.ttf",1.2
func (this *Compound) Init(width float64, height float64, fontPath string, lineSpace float64, dpi float64) error {

	this.width = width
	this.height = height

	this.lineSpace = lineSpace
	if dpi < 1 {
		dpi = 72
	}

	this.dpi = dpi

	if fontPath == "" {
		if globalFontPath == "" {
			return errors.New("Need Set FontPath")
		} else {
			this.fontPath = globalFontPath
			this.fontType = globalFontType
		}
	} else {
		this.fontPath = fontPath
		fontBytes, err := ioutil.ReadFile(this.fontPath)
		if err != nil {
			return err
		}
		this.fontType, err = truetype.Parse(fontBytes)
		if err != nil {
			return err
		}
	}

	this.obj = image.NewRGBA(image.Rect(0, 0, int(this.width), int(this.height)))
	this.Point = new(fixed.Point26_6)
	*(this.Point) = freetype.Pt(0, 0)

	draw.Draw(this.obj, this.obj.Bounds(), image.White, image.Point{}, draw.Src)
	return nil
}

func (this *Compound) splitOnSpace(x string) []string {
	var result []string
	pi := 0
	ps := false
	for i, c := range x {
		s := unicode.IsSpace(c)
		if s != ps && i > 0 {
			result = append(result, x[pi:i])
			pi = i
		}
		ps = s
	}
	result = append(result, x[pi:])
	return result
}

func (this *Compound) measureString(s string) (w, h float64) {
	a := this.fontDrawer.MeasureString(s)
	return float64(a >> 6), this.lineSpace
}

func (this *Compound) WordWrap(s string, width float64) []string {
	var result []string
	for _, line := range strings.Split(s, "\n") {
		if len(line) == 0 {
			result = append(result, "\n")
			continue
		}
		fields := this.splitOnSpace(line)

		if len(fields)%2 == 1 {
			fields = append(fields, "")
		}
		x := ""
		for i := 0; i < len(fields); i++ {
			runes := []rune(fields[i])
			for k, v := range runes {
				w, _ := this.measureString(x + string(runes[k]))
				if w > width {
					result = append(result, x)
					x = string(v)
				} else {
					x += string(v)
				}
			}
		}
		if x != "" {
			result = append(result, x)
		}
	}
	for i, line := range result {
		result[i] = strings.TrimSpace(line)
	}
	return result
}

func (this *Compound) SetX(x float64) {
	c := freetype.NewContext()
	this.Point.X = c.PointToFixed(x)
}

func (this *Compound) AddLine(fontSize float64) {
	c := freetype.NewContext()
	this.Point.Y += c.PointToFixed(fontSize * this.lineSpace)
}
func (this *Compound) AddY(c *freetype.Context, y float64) {
	this.Point.Y += c.PointToFixed(y)
}

func (this *Compound) NewContext(fontSize float64) *freetype.Context {
	this.fontDrawer = &font.Drawer{
		Face: truetype.NewFace(this.fontType, &truetype.Options{
			Size: fontSize,
		}),
	}
	c := freetype.NewContext()
	c.SetFont(this.fontType)
	c.SetFontSize(fontSize)
	c.SetDPI(this.dpi)
	c.SetClip(this.obj.Bounds())
	c.SetDst(this.obj)
	c.SetSrc(image.Black)
	return c
}

func (this *Compound) AddTitle(title string, size float64) error {
	c := this.NewContext(size)
	w, _ := this.measureString(title)
	x := (this.width - w) / 3
	if x < 0 {
		x = 0
	}

	*this.Point = freetype.Pt(int(x), 10+int(c.PointToFixed(size)>>6))
	_, err := c.DrawString(title, freetype.Pt(int(x), 10+int(c.PointToFixed(size)>>6)))
	if err != nil {
		return err
	}

	this.Point.Y += c.PointToFixed(size * this.lineSpace)
	this.Point.Y += c.PointToFixed(size * this.lineSpace)
	return nil
}

func (this *Compound) AddBody(body string, fontSize float64, width float64) error {
	c := this.NewContext(fontSize)
	this.SetX((this.width - width) / 2)
	text := this.WordWrap(body, width)
	for _, v := range text {
		_, err := c.DrawString(v, *this.Point)
		if err != nil {
			return err
		}
		this.Point.Y += c.PointToFixed(fontSize * this.lineSpace)
	}
	return nil
}

type GetImage func(imagePath string) (image.Image, error)

func GetImg(imagePath string) (image.Image, error) {
	var img image.Image
	if strings.Index(imagePath, "http") == 0 {
		resp, err := http.Get(imagePath)
		if err != nil {
			return nil, err
		}
		img, _, err = image.Decode(resp.Body)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

	} else {
		f, err := os.Open(imagePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		img, _, err = image.Decode(f)
		if err != nil {
			return nil, err
		}
	}
	return img, nil
}

// absolutePosition ture 绝对位置   false 相对位置
func (this *Compound) AddImage(imageFun GetImage, imagePath string, imageWidth uint, imageHeight uint, x int, y int, absolutePosition bool, addHeight bool) error {
	img, err := imageFun(imagePath)
	if err != nil {
		return err
	}

	img = resize.Resize(imageWidth, imageHeight, img, resize.Lanczos3)

	if absolutePosition {
		draw.Draw(this.obj, img.Bounds().Add(image.Pt(x, y)), img, img.Bounds().Min, draw.Src)
	} else {
		draw.Draw(this.obj, img.Bounds().Add(image.Pt(x, this.Point.Y.Floor()+y)), img, img.Bounds().Min, draw.Src)
	}

	if addHeight {
		c := freetype.NewContext()
		this.Point.Y += c.PointToFixed(float64(imageHeight))
	}

	return nil
}

func (this *Compound) HandleData(macroVal map[string]Data, microVal map[string]Data) map[string]Data {
	if microVal == nil {
		return macroVal
	}

	for k, v := range macroVal {
		for kk, vv := range microVal {
			if k == kk {
				v.Value = vv.Value
				macroVal[k] = v
			}
		}
	}
	return macroVal
}

type Body struct {
	Template string
	Data     map[string]Data
	Size     float64
	Width    float64
}

func (this *Compound) HandleBodies(body []Body, imageFun GetImage) error {
	for _, v := range body {
		if err := this.HandleBody(v.Template, v.Data, v.Size, v.Width, imageFun); err != nil {
			return err
		}
	}
	return nil
}

func (this *Compound) HandleBody(template string, data map[string]Data, size float64, width float64, imageFun GetImage) error {
	// 1. 正则替换文字
	re := regexp.MustCompile(`{{[a-zA-Z0-9]+}}`)
	regData := re.FindAllStringSubmatch(template, -1)
	stringData := make(map[string]Data, 0)
	for _, v := range regData {
		key := v[0]
		keyData := strings.Trim(strings.Trim(key, "{{"), "}}")
		if vv, ok := data[keyData]; !ok {
			return errors.New("Key: " + keyData + " Not Exist")
		} else {
			stringData[key] = vv
		}
	}
	for k, v := range stringData {
		if v.IsString() {
			template = strings.ReplaceAll(template, k, v.GetString())
		} else if v.IsTime() {
			template = strings.ReplaceAll(template, k, v.GetTimeString())
		}
	}

	// 2.处理图片跟换行
	re2 := regexp.MustCompile(`{{{[a-zA-Z0-9]+}}}`)
	reg2Data := re2.FindAllStringSubmatch(template, -1)
	if len(reg2Data) > 0 {
		for k, v := range reg2Data {
			index := strings.Index(template, v[0])
			if index < 0 {
				return errors.New("lines: " + v[0] + " Not Exist")
			} else {
				vv := template[0:index]

				if vv != "" {
					if err := this.AddBody(vv, size, width); err != nil {
						return err
					}
				}
				keyData := strings.Trim(strings.Trim(v[0], "{{{"), "}}}")
				if vvv, ok := data[keyData]; !ok {
					return errors.New("Key: " + keyData + " Not Exist")
				} else {
					if vvv.IsRelativeImage() {
						if err := this.AddImage(imageFun, vvv.Value, vvv.ResizeWidth, vvv.ResizeHeight, vvv.PositionX, vvv.PositionY, false, vvv.IsRise); err != nil {
							return err
						}
					} else if vvv.IsLine() {
						this.AddLine(size)
					}
				}

				template = template[index+len(v[0]):]
				if k == len(reg2Data)-1 {
					if err := this.AddBody(template, size, width); err != nil {
						return err
					}
				}
			}
		}
	} else {
		if err := this.AddBody(template, size, width); err != nil {
			return err
		}
	}
	for _, v := range data {
		if v.IsAbsoluteImage() {
			if err := this.AddImage(imageFun, v.Value, v.ResizeWidth, v.ResizeHeight, v.PositionX, v.PositionY, true, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *Compound) Save(writer io.Writer) error {
	return png.Encode(writer, this.obj)
}
