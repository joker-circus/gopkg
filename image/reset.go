package main

import (
	"image"
	"image/draw"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
)

// imaging 从目前使用上来看效率高于 resize，产生的图片对本项目无影响。

// 倍数缩放图片，percent 表示缩放百分比，10 表示为缩小至原来的 10% 大小。
func ScalingImage2(img image.Image, percent uint) image.Image {
	return ResetImage2(img, img.Bounds().Dx()*int(percent/100), img.Bounds().Dy()*int(percent/100))
}

// 缩放图片大小尺寸为 widthX * heightY，原图片大小为 img.Bounds().Dx() * img.Bounds().Dy()。
/*
	Lanczos - 用于摄影图像的高质量重采样过滤器，可产生清晰的结果。
	CatmullRom - 锐利的立方滤波器，比 Lanczos 更快更快，同时提供相似的结果。
	MitchellNetravali - 立方滤波器可产生比CatmullRom更平滑的结果，且振铃α影更少。
	Linear - 双线性重采样滤波器，产生平滑的输出。比立方过滤器更快。
	Box - 简单快速的平均滤波器适合缩小尺寸。升级时，它与最近邻类似。
	NearestNeighbor - 最快的重采样滤波器，无抗锯齿。
*/
func ResetImage2(img image.Image, widthX, heightY int) image.Image {
	return imaging.Resize(img, widthX, heightY, imaging.CatmullRom)
}

// 倍数缩放图片，percent 表示缩放百分比，10 表示为缩小至原来的 10% 大小。
func ScalingImage(img image.Image, percent uint) image.Image {
	return ResetImage(img, img.Bounds().Dx()*int(percent/100), img.Bounds().Dy()*int(percent/100))
}

// 缩放图片大小尺寸为 widthX * heightY，原图片大小为 img.Bounds().Dx() * img.Bounds().Dy()。
func ResetImage(img image.Image, widthX, heightY int) image.Image {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// 将图片等比例压缩。
	var m image.Image
	if width/height >= widthX/heightY {
		m = resize.Resize(uint(widthX), uint(height)*uint(widthX)/uint(width), img, resize.Lanczos3)
	} else {
		m = resize.Resize(uint(width*heightY/height), uint(heightY), img, resize.Lanczos3)
	}

	// 在新图片上画上压缩后的图片
	newImg := image.NewNRGBA(image.Rect(0, 0, widthX, heightY))
	if m.Bounds().Dx() > m.Bounds().Dy() {
		draw.Draw(newImg, image.Rectangle{
			Min: image.Point{Y: (heightY - m.Bounds().Dy()) / 2},
			Max: image.Point{X: widthX, Y: heightY},
		}, m, m.Bounds().Min, draw.Src)
	} else {
		draw.Draw(newImg, image.Rectangle{
			Min: image.Point{X: (widthX - m.Bounds().Dx()) / 2},
			Max: image.Point{X: widthX, Y: heightY},
		}, m, m.Bounds().Min, draw.Src)
	}
	return newImg
}
