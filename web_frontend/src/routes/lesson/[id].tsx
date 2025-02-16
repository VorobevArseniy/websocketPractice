import { useParams } from "@solidjs/router";
import { onMount } from "solid-js";
import QRGenerator from "~/components/QRCodeGenerator";
import StudentList from "~/components/StudentList";
import { formData, setFormData } from "~/lib/store";

export default function LessonPage() {
	const { id } = useParams();

	onMount(() => {
		const storedData = localStorage.getItem("formData");
		if (storedData) {
			setFormData(JSON.parse(storedData));
		}
	});

	return (
		<main class="top-0 h-screen w-full flex flex-col items-center">
			<h1>Lesson page</h1>
			<p>{id}</p>
			<p>{formData().name}</p>
			<p>{formData().subject}</p>
			<p>{formData().cabinet}</p>
			<StudentList />
			<QRGenerator text="text" />
		</main>
	);
}
