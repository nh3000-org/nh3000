package nhskin

import (
	"image/color"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type MyTheme struct{}

//var _ fyne.Theme = (*MyTheme)(nil)

var Dark = 0
var Light = 1
var Retro = 2

var Selected = 0

// var DarkBlue = color.RGBA{87, 82, 222, 1}
// var LightBlue = color.RGBA{57, 94, 169, 1}
//var Blue = color.RGBA{57, 134, 189, 1}
//var Gray = color.RGBA{186, 198, 207, 1}
//var blueColor = color.RGBA{151, 240, 173, 1}

func (m MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {

	var Blue = color.RGBA{57, 134, 189, 1}
	var Gray = color.RGBA{186, 198, 207, 1}
	//var blueColor = color.RGBA{151, 240, 173, 1}
	log.Println("skin "+strconv.Itoa(Selected)+" name ", name)
	if Selected == Dark {
		if name == "hover" {
			blueColor := color.RGBA{129, 137, 252, 1}
			return blueColor
		}
		if name == "pressed" {
			blueColor := color.RGBA{129, 137, 252, 1}
			return blueColor
		}
		if name == "selection" {
			blueColor := color.RGBA{129, 137, 252, 1}
			return blueColor
		}
		if name == "inputBackground" {
			return Gray
		}
		if name == "inputBorder" {
			return color.Black
		}
		if name == "button" {
			return Blue
		}
		if name == "foreground" {
			return color.White
		}
		if name == "background" {
			return color.Black
		}

		if name != "disabled" {

			log.Println("default ", name)
		}
	}

	if Selected == Light {
		if name == "hover" {
			blueColor := color.RGBA{129, 137, 252, 1}
			return blueColor
		}
		if name == "pressed" {
			blueColor := color.RGBA{129, 137, 252, 1}
			return blueColor
		}
		if name == "selection" {
			blueColor := color.RGBA{129, 137, 252, 1}
			return blueColor
		}
		if name == "inputBackground" {
			return color.White
		}
		if name == "inputBorder" {
			return color.Black
		}
		if name == "button" {
			return color.White
		}
		if name == "foreground" {
			return color.Black
		}
		if name == "background" {
			return color.White
		}
		if name == theme.ColorNameBackground {
			return color.Black
		}
		if name != "disabled" {

			log.Println("default ", name)
		}
	}

	if Selected == Retro {
		if name == "hover" {
			blueColor := color.RGBA{151, 240, 173, 1}
			return blueColor
		}
		if name == "selection" {
			blueColor := color.RGBA{151, 240, 173, 1}
			return blueColor
		}
		if name == "pressed" {
			greenColor := color.RGBA{151, 240, 173, 1}
			return greenColor
		}
		if name == "inputBackground" {
			return color.White
		}
		if name == "inputBorder" {
			greenColor := color.RGBA{151, 240, 173, 1}
			return greenColor
		}
		if name == "button" {
			return color.White
		}
		if name == "foreground" {
			return color.Black
		}
		if name == "background" {
			return color.White
		}
		if name == theme.ColorNameBackground {
			return color.Black
		}
		if name != "disabled" {
			log.Println("default ", name)
		}
	}
	log.Println("skin ", Selected, " ", name)
	return theme.DefaultTheme().Color(name, variant)
}
func (m MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m MyTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
func (m MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	//if name == theme.IconNameHome {
	//	fyne.NewStaticResource("myHome", homeBytes)
	//}

	return theme.DefaultTheme().Icon(name)
}
