const fetchPromise = fetch('http://localhost:8080/my-wallet/address');

fetchPromise
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error: ${response.status}`);
        }
        return response.json();
    })
    .then(json => {
            document.getElementById("address").innerHTML = json.address
            document.getElementById("pkhash").innerHTML = json.public_key_hash
            document.getElementById("pk").innerHTML = json.public_key
        });