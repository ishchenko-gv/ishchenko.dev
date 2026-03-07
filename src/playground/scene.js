import * as THREE from "three";
import { GLTFLoader } from "three/addons/loaders/GLTFLoader.js";

import { RoomEnvironment } from "three/addons/environments/RoomEnvironment.js";

const mouse = { x: 0, y: 0 };

window.addEventListener("mousemove", (event) => {
  mouse.x = (event.clientX / window.innerWidth) * 2 - 1;
  mouse.y = -(event.clientY / window.innerHeight) * 2 + 1;
});

const scene = new THREE.Scene();

const camera = new THREE.PerspectiveCamera(
  75,
  window.innerWidth / window.innerHeight,
  0.12,
  1000,
);

const light = new THREE.AmbientLight(0x112131, 0.004);
scene.add(light);
const directionalLight = new THREE.DirectionalLight(0x112131, 0.005);
directionalLight.position.set(10, 15, 10);
scene.add(directionalLight);

const loader = new GLTFLoader();

let model;

loader.load(
  "/common/twist.gltf",
  (gltf) => {
    model = gltf.scene;

    model.position.x = 3;
    model.position.y = -0.4;
    model.position.z = 0;

    model.rotation.x = 0;
    model.rotation.y = 1.6;
    model.rotation.z = 0;

    model.traverse((node) => {
      if (node.isMesh) {
        // Create or modify the material
        node.material.color.set(0x000510); // Dark charcoal/black
        node.material.metalness = 1; // Full metal
        node.material.roughness = 0.4; // Shiny/Polished

        // Boost reflections (optional)
        if (scene.environment) {
          node.material.envMapIntensity = 2.5;
        }
      }
    });

    scene.add(gltf.scene);
    console.log("Model loaded successfully");
  },
  (xhr) => {
    console.log((xhr.loaded / xhr.total) * 100 + "% loaded");
  },
  (error) => {
    console.error("An error happened", error);
  },
);

camera.position.z = 3.6;

const renderer = new THREE.WebGLRenderer({
  antialias: true,
});
renderer.setSize(window.innerWidth, window.innerHeight);

const sceneElement = document.getElementById("scene");
sceneElement.appendChild(renderer.domElement);

const pmremGenerator = new THREE.PMREMGenerator(renderer);
scene.environment = pmremGenerator.fromScene(
  new RoomEnvironment(),
  0.04,
).texture;

scene.background = new THREE.Color(0x111111);

const intensity = 0.09;

function animate(time) {
  if (model) {
    model.rotation.x += 0.001;
    model.rotation.z += 0.001;
  }
  // Calculate target position
  const targetX = mouse.x * intensity;
  const targetY = mouse.y * intensity;

  // Smoothly interpolate (Lerp) towards the target
  // 0.05 is the smoothing factor (lower = smoother/slower)
  camera.position.x += (targetX - camera.position.x) * 0.05;
  camera.position.y += (targetY - camera.position.y) * 0.05;

  // Ensure the camera always stays pointed at the center/model
  camera.lookAt(scene.position);
  renderer.render(scene, camera);
}

renderer.setAnimationLoop(animate);

window.addEventListener("resize", () => {
  // 1. Update Camera Aspect Ratio
  camera.aspect = window.innerWidth / window.innerHeight;

  // 2. Recalculate Projection Matrix (Required after changing aspect)
  camera.updateProjectionMatrix();

  // 3. Update Renderer Size
  renderer.setSize(window.innerWidth, window.innerHeight);

  // 4. Update Pixel Ratio (Optional: keeps it sharp on high-res screens)
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
});
