const latestBlockDatas = [
    {
        id: "blockRow_0",
        blockHeight: 1234,
        blockTransactions: 714,
        hash: "000000000000000000051fbcc30db5d4f09eabe02c9a461e2f9aa9c848a8c8f0",
        blockReward: "6.25",
        date: new Date('2022-09-03T03:00:48')
    }, {
        id: "blockRow_1",
        blockHeight: 342,
        blockTransactions: 1253,
        hash: "000000000000000000051fbcc30db5dasdxcqwe461e2f9aa9cszxcv848a8c8f0",
        blockReward: "2.5",
        date: new Date('2022-09-13T05:01:24')
    }
];

//LatestBlock
let tblBody = document.getElementById("latestBlockTable");
latestBlockDatas.forEach(bData => {
    let row = document.createElement("tr");
    row.setAttribute("id", bData.id);
    row.innerHTML = `
        <td><a href="#" class="blockHeight">${bData.blockHeight}</a></td>
        <td><div class="blockTransactions">${bData.blockTransactions}</div></td>
        <td><a href="#" class="hashTd">${bData.hash}</a></td>
        <td><b class="blockReward">${bData.blockReward}</b></td>
        <td><div class="dateString">${bData.date.toLocaleString()}</div></td>
    `;
    tblBody.append(row);
});
