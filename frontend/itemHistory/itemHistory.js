document.getElementById('search-form').addEventListener("submit", e => {
    e.preventDefault();
    let itHash = e.target[0].value;
    const fetchPromise = fetch(`http://localhost:8080/item/history/${itHash}`);
    fetchPromise
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error: ${response.status}`);
            }
            return response.json();
        })
        .then(json => {
            console.log(json);
            let latestMinedTransactions = json;

            let tblBody2 = document.getElementById("minedTransactions");
            latestMinedTransactions.forEach(bData => {
                let row = document.createElement("tr");
                let subRow = document.createElement("tr");
                row.setAttribute("id", bData.id);
                row.setAttribute("class", "row");
                subRow.setAttribute("class", `${bData.id}_expand expandContent`)
                row.innerHTML = `
                <td><a href="#" class="txId">${bData.txID}</a></td>
                <td><a href="#" class="txId">${bData.itemHash}</a></td>
                <td class="amountTransacted"><b class="">${bData.amount}</b><br></td>
                <td><div class="dateString">${new Date(bData.timestamp*1000).toLocaleString()}</div></td>
            `;
                subRow.innerHTML = `
                <td colspan="4">
                    <div class="expandFlex">
                        <div class="fromHolder">
                            <b>From</b>
                            <a href="#">${bData.sellerHash}</a>
                        </div>
                        <svg class="arrow" xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24">
                            <path d="M0 0h24v24H0V0z" fill="none" />
                            <path d="M8.59 16.59L13.17 12 8.59 7.41 10 6l6 6-6 6-1.41-1.41z" />
                        </svg>
                        <div class="toHolder">
                            <b>To</b>
                            <a href="#">${bData.buyerHash}</a>
                        </div>
                    </div>
                </td>
            `;
                tblBody2.append(row);
                tblBody2.append(subRow);
            });
        });
})

