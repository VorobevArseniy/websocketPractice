import { useNavigate } from "@solidjs/router";
import { Component, For, JSX, onMount } from "solid-js";
import { generateUID } from "~/lib/generateShortUID";
import { formData, setFormData } from "~/lib/store";

const inputList = ["name", "subject", "cabinet"];

const StartForm: Component<{}> = () => {
	const navigate = useNavigate();

	onMount(() => {
		if (typeof window !== "undefined") {
			const storedData = localStorage.getItem("formData");
			if (storedData) {
				setFormData(JSON.parse(storedData));
			}
		}
	});

	const setData: JSX.EventHandler<HTMLInputElement, InputEvent> = (e) => {
		setFormData({
			...formData(),
			[e.currentTarget.name]: e.currentTarget.value,
		});
	};

	const handleSubmit: JSX.EventHandler<HTMLFormElement, SubmitEvent> = (e) => {
		e.preventDefault();

		const shortUID = generateUID(6);
		localStorage.setItem("formData", JSON.stringify(formData()));

		if (typeof window !== "undefined") {
			localStorage.setItem("formData", JSON.stringify(formData()));
		}

		navigate(`/lesson/${shortUID}`);
	};
	return (
		<form
			onSubmit={handleSubmit}
			class="flex flex-col justify-center items-center gap-4 p-10">
			<For each={inputList}>
				{(field) => (
					<input
						class="border rounded-md px-2 py-1"
						placeholder={field[0].toUpperCase() + field.slice(1)}
						name={field}
						type="text"
						onInput={setData}
					/>
				)}
			</For>
			<button class="size-22 bg-gray-500 rounded-full cursor-pointer text-3xl text-white">
				+
			</button>
		</form>
	);
};

export default StartForm;
