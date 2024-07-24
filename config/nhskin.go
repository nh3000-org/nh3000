package config

import (
	"image/color"
	"log"

	//"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type MyTheme struct{}

var Dark = 0
var Light = 1
var Retro = 2

var Selected = 0

// var DarkButton = color.RGBA{162, 168, 250, 64}
var DarkButton = color.RGBA{187, 188, 201, 32}
var DarkHover = color.RGBA{187, 188, 201, 64}
var DarkPressed = color.RGBA{187, 188, 201, 220}
var DarkSelection = color.RGBA{187, 188, 201, 128}
var DarkInputBackground = color.RGBA{187, 188, 201, 32}
var DarkInputBorder = color.RGBA{187, 188, 201, 64}
var DarkSeparator = color.RGBA{187, 188, 201, 64}
var DarkShadow = color.RGBA{187, 188, 201, 64}
var DarkScrollBar = color.RGBA{187, 188, 201, 64}
var DarkFocus = color.RGBA{187, 188, 201, 64}
var DarkPlaceholder = color.RGBA{187, 188, 201, 220}
var DarkDisabled = color.RGBA{187, 188, 201, 64}
var DarkHyperlink = color.RGBA{187, 188, 201, 255}
var DarkPrimary = color.RGBA{187, 188, 201, 255}

var LightButton = color.RGBA{129, 137, 252, 250}
var LightHover = color.RGBA{129, 137, 252, 250}
var LightPressed = color.RGBA{129, 137, 252, 220}
var LightSelection = color.RGBA{129, 137, 252, 200}
var LightInputBackground = color.RGBA{129, 137, 252, 32}
var LightInputBorder = color.RGBA{129, 137, 252, 250}
var LightSeparator = color.Black
var LightShadow = color.RGBA{129, 137, 252, 64}
var LightScrollBar = color.RGBA{129, 137, 252, 250}
var LightFocus = color.RGBA{129, 137, 252, 64}
var LightPlaceholder = color.RGBA{129, 137, 252, 220}
var LightDisabled = color.RGBA{129, 137, 252, 64}
var LightHyperlink = color.RGBA{129, 137, 252, 1}
var LightPrimary = color.RGBA{129, 137, 252, 255}

var RetroButton = color.RGBA{116, 207, 103, 250}
var RetroHover = color.RGBA{116, 207, 103, 250}
var RetroPressed = color.RGBA{116, 207, 103, 220}
var RetroSelection = color.RGBA{116, 207, 103, 200}
var RetroInputBackground = color.RGBA{116, 207, 103, 32}
var RetroInputBorder = color.RGBA{116, 207, 103, 250}
var RetroSeparator = color.Black
var RetroShadow = color.RGBA{116, 207, 103, 64}
var RetroScrollBar = color.RGBA{116, 207, 103, 250}
var RetroFocus = color.RGBA{116, 207, 103, 64}
var RetroPlaceholder = color.RGBA{116, 207, 103, 255}
var RetroDisabled = color.RGBA{116, 207, 103, 64}
var RetroHyperlink = color.RGBA{116, 207, 103, 1}
var RetroPrimary = color.RGBA{116, 207, 103, 255}

func (m MyTheme) SetIcon(name fyne.ThemeIconName, variant fyne.ThemeVariant) {

}

func (m MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {

	//log.Println("skin "+strconv.Itoa(Selected)+" name ", name)
	if Selected == Dark {
		if name == "separator" {
			return DarkSeparator
		}
		if name == "shadow" {
			return DarkShadow
		}
		if name == "scrollBar" {
			return DarkScrollBar
		}
		if name == "focus" {
			return DarkFocus
		}
		if name == "placeholder" {
			return DarkPlaceholder
		}
		if name == "disabled" {
			return DarkDisabled
		}
		if name == "disabledButton" {
			return DarkDisabled
		}
		if name == "hyperlink" {
			return DarkHyperlink
		}
		if name == "primary" {
			return DarkPrimary
		}
		if name == "hover" {
			return DarkHover
		}
		if name == "pressed" {
			return DarkPressed
		}
		if name == "selection" {
			return DarkSelection
		}
		if name == "inputBackground" {
			return DarkInputBackground
		}
		if name == "inputBorder" {
			return DarkInputBorder
		}
		if name == "button" {
			return DarkButton
		}
		if name == "foreground" {
			return color.White
		}
		if name == "background" {
			return color.Black
		}
		if name == "menuBackground" {
			return color.Black
		}
		if name == "overlayBackground" {
			return color.Black
		}
		if name != "disabled" {
			log.Println("default ", name)
		}
	}

	if Selected == Light {
		if name == "separator" {
			return LightSeparator
		}
		if name == "shadow" {
			return LightShadow
		}
		if name == "scrollBar" {
			return LightScrollBar
		}
		if name == "focus" {
			return LightFocus
		}
		if name == "placeholder" {
			return LightPlaceholder
		}
		if name == "disabled" {
			return LightDisabled
		}
		if name == "disabledButton" {
			return LightDisabled
		}
		if name == "hyperlink" {
			return LightHyperlink
		}
		if name == "primary" {
			return LightPrimary
		}
		if name == "hover" {
			return LightHover
		}
		if name == "pressed" {
			return LightPressed
		}
		if name == "selection" {
			return LightSelection
		}
		if name == "inputBackground" {
			return LightInputBackground
		}
		if name == "inputBorder" {
			return LightInputBorder
		}
		if name == "button" {
			return LightButton
		}
		if name == "foreground" {
			return color.Black
		}
		if name == "background" {
			return color.White
		}
		if name == "overlayBackground" {
			return color.White
		}
		if name == "menuBackground" {
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
		if name == "separator" {
			return RetroSeparator
		}
		if name == "shadow" {
			return RetroShadow
		}
		if name == "scrollBar" {
			return RetroScrollBar
		}
		if name == "focus" {
			return RetroFocus
		}
		if name == "placeholder" {
			return RetroPlaceholder
		}
		if name == "disabled" {
			return RetroDisabled
		}
		if name == "disabledButton" {
			return RetroDisabled
		}
		if name == "hyperlink" {
			return RetroHyperlink
		}
		if name == "primary" {
			return RetroPrimary
		}
		if name == "hover" {
			return RetroHover
		}
		if name == "selection" {
			return RetroSelection
		}
		if name == "pressed" {
			return RetroPressed
		}
		if name == "inputBackground" {
			return RetroInputBackground
		}
		if name == "inputBorder" {
			return RetroInputBorder
		}
		if name == "button" {
			return RetroButton
		}
		if name == "foreground" {
			return color.Black
		}
		if name == "background" {
			return color.White
		}
		if name == "menuBackground" {
			return color.White
		}
		if name == "overlayBackground" {
			return color.White
		}
		if name == theme.ColorNameBackground {
			return color.Black
		}
		if name != "disabled" {
			log.Println("default ", name)
		}
	}
	log.Println("default ", name)
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
