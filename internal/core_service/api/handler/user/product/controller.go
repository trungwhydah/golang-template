package product

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/be/internal/core_service/api"
	productdomain "github.com/golang/be/internal/core_service/domain/product"
	"github.com/golang/be/pkg/common/httpresp"
)

type Controller struct {
	prodService productdomain.UseCaseInterface
}

func NewController(
	prodService productdomain.UseCaseInterface,
) api.Controller {
	return &Controller{
		prodService: prodService,
	}
}

func (c *Controller) RegisterRoutes(route gin.IRoutes) {
	route.GET("/products/:productId", c.GetProduct)
}

// GetProduct 	Get product by id
// @Summary 	Get product by id
// @Description Get product by id
// @Tags        product
// @Accept      json
// @Produce     json
// @Param       productId  path    string true  "Product ID"
// @Success     200  {object} httpresp.Response{data=string}
// @Failure     500  {object} httpresp.Response
// @Router      /user/products/{productId} [get].
func (c *Controller) GetProduct(g *gin.Context) {
	productID := g.Param("productId")
	if productID == "" {
		httpresp.MissingRequiredFieldError(g, "productId")

		return
	}

	curProduct, err := c.prodService.GetProduct(g, &productID)
	if err != nil {
		httpresp.InternalServerError(g)

		return
	}
	if curProduct == nil {
		httpresp.NotFound(g)

		return
	}

	res := httpresp.Response{
		Data: curProduct,
	}

	httpresp.Success(g, &res)
}
