const latestMinedTransactions = [
    {
        id: "blockRow_0",
        hash: "d2f5bd91aa389f4c7b487b0f3df89f335a65cc530d62559cdc01885da8f34681",
        amountTransacted: "6.25",
        date: new Date('2022-09-03T03:00:48'),
        from: "bc1qn4xessapgj32gxmxq6n04yfestp5j8humq0jzm",
        to: "3FntotDD41bs6t4aDxeXvyRhNB8dc8LQv4"
    }, {
        id: "blockRow_1",
        hash: "d2fasdasdas89f4c7b487b0f3df89f33asdasd30d62559cdc01885daasdasdd1",
        amountTransacted: "2.34",
        date: new Date('2022-01-11T01:00:21'),
        from: "bc1qn4xesasdasdxzczxcfestp5j8humq0jzm",
        to: "3FntotDDasdasdasdsazxcgqwyRhNB8dcqwexz"
    }
];

//latestMinedTransactions
let tblBody2 = document.getElementById("minedTransactions");
latestMinedTransactions.forEach(bData => {
    let row = document.createElement("tr");
    let subRow = document.createElement("tr");
    row.setAttribute("id", bData.id);
    row.setAttribute("class", "row");
    subRow.setAttribute("class", `${bData.id}_expand expandContent`)
    row.innerHTML = `
        <td><a href="#" class="txId">${bData.hash}</a></td>
        <td class="amountTransacted"><b class="">${bData.amountTransacted}</b><br></td>
        <td><div class="dateString">${bData.date.toLocaleString()}</div></td>
    `;
    subRow.innerHTML = `
        <td colspan="4">
            <div class="expandFlex">
                <div class="fromHolder">
                    <b>From</b>
                    <a href="#">${bData.from}</a>
                </div>
                <svg class="arrow" xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24">
                    <path d="M0 0h24v24H0V0z" fill="none" />
                    <path d="M8.59 16.59L13.17 12 8.59 7.41 10 6l6 6-6 6-1.41-1.41z" />
                </svg>
                <div class="toHolder">
                    <b>To</b>
                    <a href="#">${bData.to}</a>
                </div>
            </div>
        </td>
    `;
    tblBody2.append(row);
    tblBody2.append(subRow);
});

