-- Create function for dynamic pricing
CREATE OR REPLACE FUNCTION update_product_price()
RETURNS TRIGGER AS $func$
BEGIN
    -- Calculate demand factor based on recent orders
    WITH recent_orders AS (
        SELECT COUNT(*) as order_count
        FROM order_items oi
        JOIN orders o ON o.id = oi.order_id
        WHERE oi.product_id = NEW.id
        AND o.created_at >= NOW() - INTERVAL '7 days'
    )
    UPDATE products
    SET price = CASE
        WHEN NEW.stock <= 5 THEN base_price * 1.3
        WHEN NEW.stock <= 20 THEN base_price * 1.2
        WHEN NEW.stock <= 50 THEN base_price * 1.1
        ELSE base_price
    END * CASE
        WHEN (SELECT order_count FROM recent_orders) >= 100 THEN 1.25
        WHEN (SELECT order_count FROM recent_orders) >= 50 THEN 1.15
        WHEN (SELECT order_count FROM recent_orders) >= 20 THEN 1.1
        ELSE 1.0
    END
    WHERE id = NEW.id;

    -- Add price history record
    INSERT INTO price_histories (product_id, old_price, new_price, reason, created_at, updated_at)
    VALUES (
        NEW.id,
        OLD.price,
        NEW.price,
        CASE 
            WHEN NEW.stock <= 5 THEN 'Very low stock adjustment'
            WHEN NEW.stock <= 20 THEN 'Low stock adjustment'
            WHEN NEW.stock <= 50 THEN 'Medium stock adjustment'
            ELSE 'Regular price update'
        END,
        NOW(),
        NOW()
    );

    RETURN NEW;
END;
$func$ LANGUAGE plpgsql;

-- Create trigger for price updates
CREATE OR REPLACE TRIGGER product_price_update
AFTER UPDATE OF stock ON products
FOR EACH ROW
EXECUTE FUNCTION update_product_price();

-- Create function for inventory update
CREATE OR REPLACE FUNCTION update_inventory()
RETURNS TRIGGER AS $func$
BEGIN
    UPDATE products
    SET stock = stock - NEW.quantity
    WHERE id = NEW.product_id;
    RETURN NEW;
END;
$func$ LANGUAGE plpgsql;

-- Create trigger for inventory updates
CREATE OR REPLACE TRIGGER inventory_update
AFTER INSERT ON order_items
FOR EACH ROW
EXECUTE FUNCTION update_inventory(); 