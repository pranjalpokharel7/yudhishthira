
document.getElementById("transfer-form").addEventListener("submit", (e) => {
    e.preventDefault();
    let data = {
        item_hash: e.target[0].value + "",
        amount: Number(e.target[1].value)
    }
    console.log(data);

    postData('http://localhost:8080/transaction/coinbase', data)
        .then(data => {
            if (!data.error) {
                alert("Successful transaction")
                console.log(data);
                document.getElementById('b-hash').innerHTML = data.buyerHash;
                document.getElementById('s-hash').innerHTML = data.sellerHash;
                document.getElementById('i-hash').innerHTML = data.itemHash;
                document.getElementById('tstamp').innerHTML = data.timestamp;
                document.getElementById('txid').innerHTML = data.txID;
                document.getElementById('amount').innerHTML = data.amount;
            } else {
                alert(data.error);
            }
        });
});

async function postData(url = '', data = {}) {
    // Default options are marked with *

    const response = await fetch(url, {
        method: 'POST', // *GET, POST, PUT, DELETE, etc.
        mode: 'cors', // no-cors, *cors, same-origin
        cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
        credentials: 'same-origin', // include, *same-origin, omit
        headers: {
            'Content-Type': 'application/json'
            // 'Content-Type': 'application/x-www-form-urlencoded',
        },
        redirect: 'follow', // manual, *follow, error
        referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
        body: JSON.stringify(data) // body data type must match "Content-Type" header
    });
    return response.json(); // parses JSON response into native JavaScript objects
}



document.getElementById("clear-btn").addEventListener("click", (e) => {
    document.getElementById("transfer-form").reset()
    document.getElementById('b-hash').innerHTML = '-'
    document.getElementById('s-hash').innerHTML = '-'
    document.getElementById('i-hash').innerHTML = '-'
    document.getElementById('tstamp').innerHTML = '-'
    document.getElementById('txid').innerHTML = '-'
    document.getElementById('amount').innerHTML = '-'
})