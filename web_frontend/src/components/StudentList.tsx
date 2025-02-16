import { useParams } from "@solidjs/router";
import { Component, createSignal, For, onCleanup, onMount } from "solid-js";

const StudentList: Component = () => {
	const [users, setUsers] = createSignal<string[]>([]);
	let ws: WebSocket | null = null;
	const { id } = useParams();

	onMount(() => {
		if (typeof window === "undefined") return;

		ws = new WebSocket(`ws://localhost:8080/ws?session=${id}`);

		ws.onopen = () => {
			console.log("Подключено к WebSocket");
		};

		ws.onmessage = (event) => {
			try {
				const data = JSON.parse(event.data);
				if (Array.isArray(data)) {
					setUsers(data); // Получаем всех пользователей при подключении
				} else {
					setUsers((prev) => [...prev, data.id]); // Добавляем новых
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
