const itemHashOutput = document.getElementById("objectHash");
const verifiedHashOutput = document.getElementById("verifiedHash");

document.getElementById('hashForm').addEventListener("submit", e => {
	e.preventDefault();
	let itemID = e.target[0].value;
	let inputItemHash = e.target[1].value;
	const fetchPromise = fetch(`http://localhost:8080/item/calculate-hash/${itemID}`);
	fetchPromise
		.then(response => {
			if (!response.ok) {
				throw new Error(`HTTP error: ${response.status}`);
			}
			return response.json();
		})
		.then(json => {
			console.log(json);
			calculatedItemHash = json["item_hash"];
			itemHashOutput.innerText = calculatedItemHash;
			verifiedHashOutput.innerText = inputItemHash === calculatedItemHash;
		})
})