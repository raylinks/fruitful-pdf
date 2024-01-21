package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/olahol/go-imageupload"
)

var currentImage *imageupload.Image

func main() {

	router := gin.Default()
	//router := gin.New()
	router.LoadHTMLGlob("templates/*.html")
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	router.GET("/register", getForm)
	router.POST("/register", postForm)

	err := router.Run("0.0.0.0:" + port)
	if err != nil {
		log.Fatal(err)
	}
}

func getForm(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func postForm(c *gin.Context) {

	file, errr := imageupload.Process(c.Request, "file")
	if errr != nil {
		panic(errr)
	}

	currentImage = file

	var data [][]string

	data = append(data, []string{"Full Name", c.PostForm("fullName")})
	data = append(data, []string{"Date of birth", c.PostForm("dob")})
	data = append(data, []string{"Nationality", c.PostForm("nationality")})
	data = append(data, []string{"Street", c.PostForm("street")})
	data = append(data, []string{"City", c.PostForm("city")})
	data = append(data, []string{"State", c.PostForm("state")})
	data = append(data, []string{"Country", c.PostForm("country")})
	data = append(data, []string{"Postal Code", c.PostForm("postalCode")})
	data = append(data, []string{"Phone Number", c.PostForm("phoneNumber")})
	data = append(data, []string{"Email", c.PostForm("email")})
	data = append(data, []string{"Document Type", c.PostForm("documentType")})
	data = append(data, []string{"Document Number", c.PostForm("documentNumber")})
	data = append(data, []string{"Issuing Country", c.PostForm("issuingCountry")})
	data = append(data, []string{"Expiry Date", c.PostForm("expiryDate")})
	data = append(data, []string{"file", c.PostForm("file")})

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(20, 10, 20)
	buildHeading(m)
	buildFruitList(m, data)

	timestamp := time.Now().Format("20060102150405")

	fileName := "kyc-" + timestamp + ".pdf"
	filePath := "pdfs/" + fileName
	err := m.OutputFileAndClose(filePath)
	if err != nil {
		fmt.Println("could not save pdf", err)
		os.Exit(1)
	}

	fmt.Println("PDF saved successfully")
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Writer.Header().Add("Content-Type", "application/pdf")
	//c.Writer(c.currentImage)
	c.File(filePath)
}

func buildFruitList(m pdf.Maroto, contents [][]string) {
	tableHeadings := []string{"Name", "Details"}

	fmt.Println(contents)

	lightPurpleColor := getLightPurpleColor()

	m.SetBackgroundColor(getTealColor())
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("KYC(Know Your Customers)", props.Text{
				Top:    2,
				Size:   13,
				Color:  color.NewWhite(),
				Family: consts.Courier,
				Style:  consts.Bold,
				Align:  consts.Center,
			})
		})
	})

	m.SetBackgroundColor(color.NewWhite())
	m.TableList(tableHeadings, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{3, 7, 2},
		},
		ContentProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{3, 7, 2},
		},
		Align:                consts.Left,
		HeaderContentSpace:   1,
		Line:                 false,
		AlternatedBackground: &lightPurpleColor,
	})
}

func buildHeading(m pdf.Maroto) {
	m.RegisterHeader(func() {
		m.Row(50, func() {
			m.Col(12, func() {
				err := m.FileImage("images/ben.jpeg", props.Rect{
					Center:  true,
					Percent: 75,
				})

				if err != nil {
					fmt.Println("image file not uploaded", err)
				}
			})
		})
	})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Please fill up the KYC form...", props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Center,
				Color: getDarkPurpleColor(),
			})
		})
	})
}

func getDarkPurpleColor() color.Color {
	return color.Color{
		Red:   88,
		Green: 80,
		Blue:  99,
	}
}
func getTealColor() color.Color {
	return color.Color{
		Red:   3,
		Green: 166,
		Blue:  166,
	}
}

func getLightPurpleColor() color.Color {
	return color.Color{
		Red:   210,
		Green: 200,
		Blue:  230,
	}
}
