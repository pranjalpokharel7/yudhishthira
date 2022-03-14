const signedTokenOutput = document.getElementById("signedToken");

document.getElementById('token-form').addEventListener("submit", e => {
	e.preventDefault();
	let token = e.target[0].value;
	const fetchPromise = fetch(`http://localhost:8080/token/sign/${token}`);
	fetchPromise
		.then(response => {
			if (!response.ok) {
				throw new Error(`HTTP error: ${response.status}`);
			}
			return response.json();
		})
		.then(json => {
			console.log(json);
			signedTokenOutput.innerText = json["signed_token"];
			
		})
})