const { tool } = require("@langchain/core/tools")
const { z } = require("zod")
const axios = require("axios")

const searchProduct = tool(async ({ query, token }) => {

    console.log("searchProduct called with data:", { query, token })

    const response = await axios.get(`${process.env.PRODUCT_SERVICE_URL}/api/product/get?q=${query}`, {
        headers: {
            Authorization: `Bearer ${token}`
        }
    })

    return JSON.stringify(response.data)

}, {

    name: "searchProduct",
    description: "Search for products based on a query",
    schema: z.object({
        query: z.string().describe("The search query for products")
    })
})


const addProductToCart = tool(async ({ productId, quantity = 1, token }) => {


    const response = await axios.post(`${process.env.CART_SERVICE_URL}/api/cart/item`, {
        productId,
        quantity
    }, {
        headers: {
            Authorization: `Bearer ${token}`
        }
    })

    return `Added product with id ${productId} (qty: ${quantity}) to cart`
 

}, {
    name: "addProductToCart",
    description: "Add a product to the shopping cart",
    schema: z.object({
        productId: z.string().describe("The id of the product to add to the cart"),
        quantity: z.number().describe("The quantity of the product to add to the cart").default(1),
    })
})


module.exports = { searchProduct, addProductToCart }