package interfigures

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

// круг
type Circle struct {
	Radius float64 `json:"radius"`
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// квадрат
type Square struct {
	Side float64 `json:"side"`
}

func (s Square) Area() float64 {
	return s.Side * s.Side
}

func (s Square) Perimeter() float64 {
	return s.Side * 4
}

// прямоугольник
type Rectangle struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// принт
func PrintShape(s Shape) {
	switch v := s.(type) {
	case Circle:
		fmt.Printf("Круг: площадь = %.2f, периметр = %.2f\n", v.Area(), v.Perimeter())
	case Square:
		fmt.Printf("Квадрат: площадь = %.2f, периметр = %.2f\n", v.Area(), v.Perimeter())
	case Rectangle:
		fmt.Printf("Прямоугольник: площадь = %.2f, периметр = %.2f\n", v.Area(), v.Perimeter())
	default:
		fmt.Printf("Неизвестная фигура: площадь = %.2f, периметр = %.2f\n", s.Area(), s.Perimeter())
	}
}

type ShapeWrapper struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func CreateShape(wrapper ShapeWrapper) (Shape, error) {
	switch wrapper.Type {
	case "круг":
		var c Circle
		err := json.Unmarshal(wrapper.Data, &c)
		return c, err
	case "квадрат":
		var s Square
		err := json.Unmarshal(wrapper.Data, &s)
		return s, err
	case "прямоугольник":
		var r Rectangle
		err := json.Unmarshal(wrapper.Data, &r)
		return r, err
	default:
		return nil, fmt.Errorf("неизвестный тип: %s", wrapper.Type)
	}
}

func ReadFigures() ([]Shape, error) {
	data, err := os.ReadFile("interfigures/figures.json")
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл: %w", err)
	}

	var wrappers []ShapeWrapper
	err = json.Unmarshal(data, &wrappers)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	var shapes []Shape
	for _, w := range wrappers {
		shape, err := CreateShape(w)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания фигуры %s: %w", w.Type, err)
		}
		shapes = append(shapes, shape)
	}

	return shapes, nil
}
