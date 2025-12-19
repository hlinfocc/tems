package captcha

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math/rand"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// CaptchaGenerator 验证码生成器
type CaptchaGenerator struct {
	Width    int
	Height   int
	CodeLen  int
	FontFace font.Face
}

// NewCaptchaGenerator 创建验证码生成器
func NewCaptchaGenerator(width, height, codeLen int) *CaptchaGenerator {
	return &CaptchaGenerator{
		Width:    width,
		Height:   height,
		CodeLen:  codeLen,
		FontFace: basicfont.Face7x13, // 使用内置字体
	}
}

// Generate 生成验证码图片和验证码文本
func (cg *CaptchaGenerator) Generate() (image.Image, string, error) {
	code := cg.generateRandomCode()
	img := cg.createCaptchaImage(code)
	return img, code, nil
}

// 生成随机字符串（数字+大写字母）
func (cg *CaptchaGenerator) generateRandomCode() string {
	const charset = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZ" // 避免易混淆字符
	code := make([]byte, cg.CodeLen)

	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// 创建验证码图片
func (cg *CaptchaGenerator) createCaptchaImage(code string) image.Image {
	// 创建空白画布
	img := image.NewRGBA(image.Rect(0, 0, cg.Width, cg.Height))

	// 填充背景色
	bgColor := color.RGBA{240, 240, 245, 255} // 浅灰色背景
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 绘制干扰元素
	cg.drawInterference(img)

	// 绘制验证码字符
	cg.drawCode(img, code)

	return img
}

// 绘制干扰元素（线和点）
func (cg *CaptchaGenerator) drawInterference(img *image.RGBA) {
	// 绘制干扰线
	for i := 0; i < 3; i++ {
		cg.drawRandomLine(img)
	}

	// 添加噪点
	cg.addNoise(img)
}

// 绘制随机干扰线
func (cg *CaptchaGenerator) drawRandomLine(img *image.RGBA) {
	// 随机颜色
	lineColor := color.RGBA{
		uint8(rand.Intn(100) + 50),
		uint8(rand.Intn(100) + 50),
		uint8(rand.Intn(100) + 50),
		180, // 半透明
	}

	// 随机起点终点
	x1 := rand.Intn(cg.Width / 4)
	y1 := rand.Intn(cg.Height)
	x2 := cg.Width*3/4 + rand.Intn(cg.Width/4)
	y2 := rand.Intn(cg.Height)

	// 绘制直线（简单实现）
	cg.drawLine(img, x1, y1, x2, y2, lineColor)
}

// 绘制直线 (Bresenham 算法)
func (cg *CaptchaGenerator) drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx, sy := 1, 1
	if x1 > x2 {
		sx = -1
	}
	if y1 > y2 {
		sy = -1
	}
	err := dx - dy

	for {
		img.Set(x1, y1, col)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// 添加噪点
func (cg *CaptchaGenerator) addNoise(img *image.RGBA) {
	pointCount := cg.Width * cg.Height / 10 // 10%像素点作为噪点
	for i := 0; i < pointCount; i++ {
		x := rand.Intn(cg.Width)
		y := rand.Intn(cg.Height)
		rnd := rand.Intn(100)

		// 随机颜色噪点
		if rnd < 70 { // 70%概率为深色点
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		} else { // 30%概率为随机颜色点
			img.Set(x, y, color.RGBA{
				uint8(rand.Intn(150)),
				uint8(rand.Intn(150)),
				uint8(rand.Intn(150)),
				255,
			})
		}
	}
}

// 绘制验证码字符
func (cg *CaptchaGenerator) drawCode(img *image.RGBA, code string) {
	charWidth := cg.Width / len(code)
	fontHeight := cg.FontFace.Metrics().Height.Ceil()

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: cg.FontFace,
	}

	for i, char := range code {
		// 随机字符颜色
		charColor := color.RGBA{
			uint8(rand.Intn(100) + 30), // R: 30-130
			uint8(rand.Intn(100) + 30), // G: 30-130
			uint8(rand.Intn(100) + 30), // B: 30-130
			255,
		}

		d.Src = image.NewUniform(charColor)

		// 字符位置（居中）
		x := i*charWidth + (charWidth-cg.FontFace.Metrics().Ascent.Ceil())/2
		y := (cg.Height+fontHeight)/2 - 25 // 垂直居中

		// 添加随机偏移
		x += rand.Intn(5) - 2
		y += rand.Intn(5) - 2

		// 绘制字符
		// d.Dot = fixed.P(x, y)
		// d.DrawString(string(char))
		// 绘制缩放后的字符
		cg.drawScaledChar(img, rune(char), x, y, charColor, 3)
	}
}

// 绘制缩放后的字符
func (cg *CaptchaGenerator) drawScaledChar(img *image.RGBA, char rune, x, y int, col color.Color, scale float64) {
	// 使用基本字体
	face := basicfont.Face7x13

	// 获取字符的位图
	dot := fixed.P(0, 0)
	dr, mask, maskp, _, ok := face.Glyph(dot, rune(char))
	if !ok {
		return
	}

	// 绘制每个像素，应用缩放
	for py := dr.Min.Y; py < dr.Max.Y; py++ {
		for px := dr.Min.X; px < dr.Max.X; px++ {
			// 获取原始像素值
			_, _, _, a := mask.At(px-dr.Min.X+maskp.X, py-dr.Min.Y+maskp.Y).RGBA()

			// 如果像素可见
			if a > 0 {
				// 应用缩放
				for sy := 0; sy < int(scale); sy++ {
					for sx := 0; sx < int(scale); sx++ {
						targetX := x + (px-dr.Min.X)*int(scale) + sx
						targetY := y + (py-dr.Min.Y)*int(scale) + sy

						// 确保在图像范围内
						if targetX >= 0 && targetX < cg.Width && targetY >= 0 && targetY < cg.Height {
							img.Set(targetX, targetY, col)
						}
					}
				}
			}
		}
	}
}

// ImageToJPEG 将图片转换为JPEG字节流
func (cg *CaptchaGenerator) ImageToJPEG(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 85})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateBase64 生成Base64编码的验证码图片和验证码文本
func GenerateBase64() (string, string, error) {
	captchaGen := NewCaptchaGenerator(145, 50, 5)
	img, code, err := captchaGen.Generate()
	if err != nil {
		return "", "", err
	}

	// 将图片转换为Base64
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 85})
	if err != nil {
		return "", "", err
	}

	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/jpeg;base64," + base64Str, code, nil
}
