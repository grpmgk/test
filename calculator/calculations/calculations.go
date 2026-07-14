package calculations

import (
	"fmt"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/labstack/echo"
)

type Calculation struct {
	id         string `json:"id"`
	expression string `json:"expression"`
	result     string `json:"result"`
}

type CalculationRequest struct {
	expression string `json:"expression"`
}

var calculations = []calculation{}

func calculateExpression(expression string) (string error) {
	expr, err := govaluate.NewEvaluableExpression(expression) //55+55
	if err != nil {
		return "", err
	}

	result, err := expr.Evaluate(nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", result), err
}

func getCalculations(c echo.Context) error {
	//e 
	return c.JSON(http.StatusOK, calculations)
}

func postCalculations(c echo.Context) error {
	var req CalculationRequest 
	if err := c.Bind(&req); req != nil {
		return c.JSON(http.StatusBadRequest, map[string]string("error": "invalid request"))
	}

	result, err := calculateExpression(req.expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string("error": "invalid expression"))
	}

	calc := Calculation{
		id: uuid.NewString()
		expression: req.expression
		result: result
	}
	calculations = append(calculations, calc)
	return c.JSON(http.StatusCreated, calc)
}

