import { Component, onMount } from "solid-js";
import QRCode from "qrcode";

interface QRGeneratorProps {
	text: string;
}

const QRGenerator: Component<QRGeneratorProps> = (props) => {
	let canvasRef: HTMLCanvasElement | undefined;

	onMount(() => {
		if (canvasRef) {
			QRCode.toCanvas(canvasRef, props.text, { width: 200 })
				.then(() => console.log("QR Code generated"))
				.catch((err) => console.error("QR Code error:", err));
		}
	});

	return (
		<div class="flex flex-col items-center justify-center p-4">
			<h1 class="text-xl font-bold mb-4">QR Code Generator</h1>
			<canvas ref={canvasRef} class="border border-gray-300" />
		</div>
	);
};

export default QRGenerator;
