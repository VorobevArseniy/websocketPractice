import QRGenerator from "~/components/QRCodeGenerator";
import StartForm from "~/components/StartForm";

export default function Home() {
	return (
		<main class="top-0 h-screen w-full mx-auto flex flex-col justify-center items-center text-gray-700 p-4">
			<StartForm />
			<QRGenerator text="text" />
		</main>
	);
}
