import { useParams } from "@solidjs/router";
import { Component, createSignal, For, onCleanup, onMount } from "solid-js";
import { formData, setSessionID } from "~/lib/store";

interface studentInterface {
	ID: string;
	Name: string;
}

const StudentList: Component = () => {
	const [users, setUsers] = createSignal<studentInterface[]>([]);
	let ws: WebSocket | null = null;
	const { id } = useParams();

	onMount(() => {
		if (typeof window === "undefined") return;

		setSessionID(id);

		ws = new WebSocket(
			`ws://localhost:8080/ws?session=${id}&teacher=${formData().name}`,
		);

		ws.onopen = () => {
			console.log("Подключено к WebSocket");
		};

		ws.onmessage = (event) => {
			try {
				const data = JSON.parse(event.data);
				console.log(data);
				if (Array.isArray(data)) {
					setUsers(data); // Получаем всех пользователей при подключении
				} else {
					setUsers((prev) => [...prev, data]); // Добавляем новых
				}
			} catch (err) {
				console.error(err);
			}
		};

		ws.onclose = (event) => {
			console.log("WebSocket закрыт:", event.code, event.reason);
		};

		ws.onerror = (error) => {
			console.error("Ошибка WebSocket:", error);
		};

		onCleanup(() => {
			console.log("closed ws");
			if (ws) {
				ws.close();
			}
		});
	});

	return (
		<ul>
			<For each={users()}>{(user) => <li>{user}</li>}</For>
		</ul>
	);
};

export default StudentList;
