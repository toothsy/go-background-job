import { useEffect, useState } from "react";
import "./App.css";
import axios from "axios";
function App() {
	const [text, setText] = useState<string>("");
	const [isLoading, setIsLoading] = useState<boolean>(true);

	useEffect(() => {
		const fetchData = async () => {
			try {
				let body = await axios.get("http://localhost:8080/");
				console.log(body);
				setText(body.data);
			} catch (e) {
				console.log(e);
				throw e;
			} finally {
				setIsLoading(false);
			}
		};

		fetchData();
	}, []);
	return (
		<>
			<div className="flex flex-row justify-center">Hallo</div>
			{isLoading ? <div>Loading...</div> : <div>{text}</div>}
		</>
	);
}

export default App;
