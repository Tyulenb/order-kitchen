let order = {
    quantityBurger: 0,
    quantityFries: 0,
    quantityCola: 0,
};

function adjustQuantity(productKey, change) {
    order[productKey] = Math.max(0, order[productKey] + change);
    document.getElementById('display' + productKey.charAt(0).toUpperCase() + productKey.slice(1)).textContent = order[productKey];
}

async function makeOrder() {
    try {
        console.log(JSON.stringify(order))

        const response = await fetch('/order', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(order)
        });

        if (!response.ok) {
            throw new Error('Network response was not ok: ' + errorText);
        }

    } catch (error) {
        alert('Error placing the order. Please try again.');
    }
}
