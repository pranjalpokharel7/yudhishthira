
const fetchPromise = fetch('http://localhost:8080/block/last/10');

fetchPromise
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error: ${response.status}`);
        }
        return response.json();
    })
    .then(json => {
        let latestBlockDatas = json;

        let tblBody = document.getElementById("latestBlockTable");
        latestBlockDatas.forEach(bData => {
            let row = document.createElement("tr");
            row.setAttribute("id", bData.id);
            row.innerHTML = `
                <td><a href="#" class="blockHeight">${bData.height}</a></td>
                <td><div class="blockTransactions">${bData.nonce}</div></td>
                <td><a href="#" class="hashTd">${bData.block_hash}</a></td>
                <td><b class="blockReward">${bData.difficulty}</b></td>
                <td><div class="dateString">${new Date(bData.timestamp*1000).toLocaleString()}</div></td>
            `;
            tblBody.append(row);
        });
    });


// const URL = 'http://localhost:8080/block/last/3';

// async function getData() {
//     const response = await fetch(this.URI);
//     const data = await response.json();
//     return data;
// }

// console.log(getData())

//LatestBlock

{/* <td><div class="dateString">${bData.date.toLocaleString()}</div></td> */ }
