-- name: ListOrders :many
SELECT * FROM orders ORDER BY order_date DESC;


-- name: ListOrdersMonthly :many
SELECT
    m.month::date AS month,
    COALESCE(SUM(o.total), 0)::numeric AS total,
    COUNT(o.order_date) AS total_orders
FROM generate_series(
    make_date(sqlc.arg(start_year)::int, 1, 1),
    date_trunc('month', CURRENT_DATE)::date,
    interval '1 month'
) AS m(month)
LEFT JOIN orders o
    ON o.order_date >= m.month AND o.order_date < m.month + interval '1 month'
GROUP BY m.month
ORDER BY m.month;

-- name: ListFavRestaurants :many
SELECT
  COALESCE(SUM(total), 0)::numeric AS total,
  COUNT(*) AS order_count,
  restaurant_name
FROM orders
WHERE order_date >= sqlc.arg(start_at)::timestamptz
  AND order_date < sqlc.arg(end_at)::timestamptz
GROUP BY restaurant_name
ORDER BY total DESC, restaurant_name
LIMIT 10;

-- name: GetLastOrderTime :one
SELECT order_date FROM orders ORDER BY order_date DESC LIMIT 1;

-- name: ListOrderYears :many
SELECT DISTINCT EXTRACT(YEAR FROM order_date)::int AS year
FROM orders
ORDER BY year DESC;

-- name: GetPeriodSummary :one
SELECT
  COALESCE(SUM(total), 0)::numeric AS total,
  COUNT(*) AS order_count
FROM orders
WHERE order_date >= sqlc.arg(start_at)::timestamptz
  AND order_date < sqlc.arg(end_at)::timestamptz;

-- name: ListOrdersBetween :many
SELECT order_date, total
FROM orders
WHERE order_date >= sqlc.arg(start_at)::timestamptz
  AND order_date < sqlc.arg(end_at)::timestamptz
ORDER BY order_date;

-- name: GetFirstOrderTime :one
SELECT order_date FROM orders ORDER BY order_date ASC LIMIT 1;
