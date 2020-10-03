package context

import (
	"github.com/labstack/echo/v4"
)

/**
 * CustomContext
 */
type CustomContext struct {
	echo.Context
}

/**
 * Bar
 */
func (c *CustomContext) Bar() string {
	return "bar"
}
