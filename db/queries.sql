-- name: ListOrders :many
SELECT * FROM orders ORDER BY order_date DESC;
