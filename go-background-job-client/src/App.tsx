import { useEffect, useState } from "react";
import "./App.css";
import axios from "axios";
import { useQuery } from "react-query";
function App() {
	const [text, setText] = useState<string>("");
	const { isLoading, error, data } = useQuery({
		queryKey: ["repoData"],
		queryFn: () => axios.get("http://localhost:8080/"),
	});
	useEffect(() => {
		if (error) console.log(error);
	}, [error]);
	useEffect(() => {
		if (!isLoading) setText(data?.data);
	}, [isLoading]);
	return (
		<>
			<div className="flex flex-row justify-center">Hallo</div>
			{isLoading ? <div>Loading...</div> : <div>{text}</div>}
		</>
	);
}

export default App;
