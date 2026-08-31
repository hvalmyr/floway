import * as THREE from "three";

/**
 * Procedural "minimalist flowery branch" — ported from the Claude Design
 * canvas (project 61bbe33f, file flowery-branch.js). Materials already use
 * the site's exact 4-color palette (bark/leaf ≈ ink, accent = primary), so
 * it reads as on-brand without extra styling.
 *
 * `seed` drives every random choice below — pass a fixed value for a
 * reproducible tree (2024 is the original design's default) or a random
 * one for a fresh composition. Not every seed grows something pleasant;
 * pair with `isReasonableBranch` and retry on a new seed when it fails.
 */
export function buildFloweryBranch(seed = 2024): THREE.Group {
  const bark = new THREE.MeshStandardMaterial({
    color: 0x41342a,
    roughness: 0.85,
    metalness: 0.05,
  });
  bark.name = "bark";
  const leaf = new THREE.MeshStandardMaterial({
    color: 0x3a2f24,
    roughness: 0.55,
    metalness: 0.0,
    side: THREE.DoubleSide,
  });
  leaf.name = "leaf";
  const petal = new THREE.MeshStandardMaterial({ color: 0xd9e8f0, roughness: 0.4, metalness: 0.0 });
  petal.name = "petal";
  const accent = new THREE.MeshStandardMaterial({
    color: 0x82b1cc,
    roughness: 0.35,
    metalness: 0.1,
  });
  accent.name = "accent";

  const branch = new THREE.Group();
  branch.name = "flowery_branch";

  // Seeded RNG for reproducible-but-organic variation.
  function mulberry32(seed: number) {
    return function () {
      seed |= 0;
      seed = (seed + 0x6d2b79f5) | 0;
      let t = Math.imul(seed ^ (seed >>> 15), 1 | seed);
      t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
      return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
    };
  }
  const rng = mulberry32(seed);

  function v(x: number, y: number, z: number) {
    return new THREE.Vector3(x, y, z);
  }
  function dirQuat(dir: THREE.Vector3) {
    return new THREE.Quaternion().setFromUnitVectors(v(0, 1, 0), dir.clone().normalize());
  }

  function quatFromAxes(yAxis: THREE.Vector3, refAxis: THREE.Vector3) {
    const y = yAxis.clone().normalize();
    let x = new THREE.Vector3().crossVectors(refAxis, y);
    if (x.lengthSq() < 1e-6) x = new THREE.Vector3().crossVectors(v(1, 0, 0), y);
    x.normalize();
    const z = new THREE.Vector3().crossVectors(x, y).normalize();
    return new THREE.Quaternion().setFromRotationMatrix(new THREE.Matrix4().makeBasis(x, y, z));
  }

  function randomPerp(dir: THREE.Vector3) {
    let ref = v(0.4, 1, 0.1);
    if (Math.abs(dir.clone().normalize().dot(ref.clone().normalize())) > 0.9) ref = v(1, 0, 0.3);
    return new THREE.Vector3().crossVectors(dir, ref).normalize();
  }

  let meshCount = 0;
  function addSegment(
    group: THREE.Group,
    p0: THREE.Vector3,
    p1: THREE.Vector3,
    r0: number,
    r1: number,
    mat: THREE.Material,
    name: string,
    radialSegments: number,
  ) {
    const dir = new THREE.Vector3().subVectors(p1, p0);
    const len = dir.length();
    if (len < 1e-6) return;
    const geo = new THREE.CylinderGeometry(r1, r0, len, radialSegments, 1, false);
    geo.translate(0, len / 2, 0);
    const mesh = new THREE.Mesh(geo, mat);
    mesh.name = `${name}_${meshCount++}`;
    mesh.position.copy(p0);
    mesh.quaternion.copy(dirQuat(dir));
    group.add(mesh);
  }

  function addChain(
    group: THREE.Group,
    points: THREE.Vector3[],
    radii: number[],
    mat: THREE.Material,
    name: string,
    radialSegments = 10,
  ) {
    for (let i = 0; i < points.length - 1; i++) {
      addSegment(
        group,
        points[i],
        points[i + 1],
        radii[i],
        radii[i + 1],
        mat,
        name,
        radialSegments,
      );
    }
  }

  function leafGeometry(length: number, width: number, thickness: number) {
    const shape = new THREE.Shape();
    shape.moveTo(0, 0);
    shape.quadraticCurveTo(width * 0.55, length * 0.32, width * 0.06, length * 0.82);
    shape.quadraticCurveTo(0, length * 0.95, 0, length);
    shape.quadraticCurveTo(0, length * 0.95, -width * 0.06, length * 0.82);
    shape.quadraticCurveTo(-width * 0.55, length * 0.32, 0, 0);
    const geo = new THREE.ExtrudeGeometry(shape, {
      depth: thickness,
      bevelEnabled: false,
      curveSegments: 10,
    });
    geo.translate(0, 0, -thickness / 2);
    return geo;
  }

  function addLeaf(
    group: THREE.Group,
    base: THREE.Vector3,
    outwardDir: THREE.Vector3,
    size: number,
    name: string,
    droopDeg = 28,
  ) {
    const geo = leafGeometry(size, size * 0.4, size * 0.025);
    const mesh = new THREE.Mesh(geo, leaf);
    mesh.name = `${name}_${meshCount++}`;
    const baseQuat = quatFromAxes(outwardDir, v(0, 1, 0.3));
    const droop = new THREE.Quaternion().setFromAxisAngle(
      v(1, 0, 0),
      -THREE.MathUtils.degToRad(droopDeg),
    );
    mesh.quaternion.copy(baseQuat.multiply(droop));
    mesh.position.copy(base);
    group.add(mesh);
  }

  // Tiny bloomed flower: radial petals + center, facing along axisDir.
  function addFlower(
    group: THREE.Group,
    tip: THREE.Vector3,
    axisDir: THREE.Vector3,
    size: number,
    petalCount: number,
    name: string,
  ) {
    const axis = axisDir.clone().normalize();
    let ref = v(0.3, 1, 0);
    if (Math.abs(axis.dot(ref)) > 0.95) ref = v(1, 0, 0.2);
    const u = new THREE.Vector3().crossVectors(ref, axis).normalize();
    const w = new THREE.Vector3().crossVectors(axis, u).normalize();

    const petalLen = size * 1.3;
    const petalWidth = size * 0.16;
    const petalThick = size * 0.07;
    for (let i = 0; i < petalCount; i++) {
      const ang = (i / petalCount) * Math.PI * 2;
      const radial = new THREE.Vector3()
        .addScaledVector(u, Math.cos(ang))
        .addScaledVector(w, Math.sin(ang))
        .normalize();
      const petalDir = new THREE.Vector3()
        .addScaledVector(radial, 0.92)
        .addScaledVector(axis, 0.34)
        .normalize();
      const geo = new THREE.SphereGeometry(1, 12, 10);
      geo.scale(petalWidth, petalLen, petalThick);
      geo.translate(0, petalLen, 0);
      const mesh = new THREE.Mesh(geo, petal);
      mesh.name = `${name}_petal_${i}`;
      mesh.quaternion.copy(dirQuat(petalDir));
      mesh.position.copy(tip);
      group.add(mesh);
    }
    const centerMesh = new THREE.Mesh(new THREE.SphereGeometry(size * 0.14, 14, 10), accent);
    centerMesh.name = `${name}_center`;
    centerMesh.position.copy(tip).addScaledVector(axis, size * 0.16);
    group.add(centerMesh);
  }

  function addBud(
    group: THREE.Group,
    tip: THREE.Vector3,
    axisDir: THREE.Vector3,
    size: number,
    name: string,
  ) {
    const axis = axisDir.clone().normalize();
    const geo = new THREE.SphereGeometry(1, 12, 10);
    geo.scale(size * 0.34, size * 0.62, size * 0.34);
    geo.translate(0, size * 0.5, 0);
    const mesh = new THREE.Mesh(geo, petal);
    mesh.name = `${name}_bud`;
    mesh.quaternion.copy(dirQuat(axis));
    mesh.position.copy(tip);
    group.add(mesh);

    const cap = new THREE.Mesh(new THREE.SphereGeometry(size * 0.1, 10, 8), accent);
    cap.name = `${name}_bud_cap`;
    cap.position.copy(tip).addScaledVector(axis, size * 0.95);
    group.add(cap);
  }

  // Recursive branch growth.
  let tipCount = 0;
  function growBranch(
    start: THREE.Vector3,
    dir: THREE.Vector3,
    length: number,
    radius: number,
    depth: number,
    name: string,
  ) {
    const segs = 4 + Math.floor(rng() * 2);
    const pts = [start.clone()];
    const radii = [radius];
    let curDir = dir.clone().normalize();
    let curPos = start.clone();
    const segLen = length / segs;

    for (let i = 0; i < segs; i++) {
      const perturb = v((rng() - 0.5) * 0.4, (rng() - 0.15) * 0.3, (rng() - 0.5) * 0.4);
      curDir = curDir.clone().add(perturb).normalize();
      curPos = curPos.clone().add(curDir.clone().multiplyScalar(segLen));
      pts.push(curPos.clone());
      radii.push(radius * (1 - ((i + 1) / segs) * 0.42));
    }

    const radialSegs = radius > 0.004 ? 14 : radius > 0.002 ? 10 : 8;
    addChain(branch, pts, radii, bark, name, radialSegs);

    const endPos = pts[pts.length - 1];
    const endDir = curDir;
    const endRadius = radii[radii.length - 1];

    if (depth <= 2 && rng() < 0.4) {
      const li = 1 + Math.floor(rng() * (pts.length - 2));
      const outward = randomPerp(curDir)
        .multiplyScalar(0.7)
        .add(curDir.clone().multiplyScalar(0.2));
      addLeaf(branch, pts[li], outward, 0.028 + rng() * 0.014, `${name}_leaf`, 25 + rng() * 20);
    }

    if (depth > 0 && endRadius > 0.0003) {
      const branchChance = depth >= 2 ? 0.97 : 0.9;
      const nChildren = rng() < branchChance ? 2 : 1;
      for (let c = 0; c < nChildren; c++) {
        const spread = 0.45 + rng() * 0.55;
        const axis = randomPerp(endDir);
        const childDir = endDir.clone().applyAxisAngle(axis, spread * (c % 2 === 0 ? 1 : -1));
        childDir.applyAxisAngle(endDir, rng() * Math.PI * 2 * 0.6);
        const childLen = length * (0.6 + rng() * 0.18);
        const childRadius = endRadius * (0.7 + rng() * 0.18);
        growBranch(endPos, childDir, childLen, childRadius, depth - 1, `${name}_${c}`);
      }
      // Occasional small blossom mid-fork for a scattered-bloom look.
      if (rng() < 0.3)
        addFlower(branch, endPos, endDir, 0.01 + rng() * 0.006, 5, `bloom_${tipCount++}`);
    } else {
      if (rng() < 0.62)
        addFlower(branch, endPos, endDir, 0.011 + rng() * 0.007, 5, `bloom_${tipCount++}`);
      else addBud(branch, endPos, endDir, 0.02 + rng() * 0.008, `bud_${tipCount++}`);
    }
  }

  // Asymmetric leaning trunk, ikebana-style single composition.
  growBranch(v(0, 0, 0), v(0.12, 1, 0.05), 0.26, 0.0125, 6, "trunk");

  return branch;
}

/**
 * Sampled 30 random seeds while building this — the growth probabilities
 * are high enough (0.9–0.97) that outright bare/broken trees are rare, but
 * not impossible. Cheap sanity checks to reject the rare bad roll: enough
 * mesh detail, at least a handful of visible blooms/buds (a "flowery"
 * branch with none looks like a bug), and not a wildly skewed silhouette.
 */
export function isReasonableBranch(branch: THREE.Group): boolean {
  let meshCount = 0;
  let bloomCount = 0;
  branch.traverse((o) => {
    if (!(o as THREE.Mesh).isMesh) return;
    meshCount++;
    if (o.name.includes("_center") || o.name.includes("_bud_cap")) bloomCount++;
  });
  if (meshCount < 150 || meshCount > 900) return false;
  if (bloomCount < 5) return false;

  const box = new THREE.Box3().setFromObject(branch);
  const size = box.getSize(new THREE.Vector3());
  const horizontal = Math.max(size.x, size.z) || 0.0001;
  const ratio = size.y / horizontal;
  if (ratio > 4.5 || ratio < 0.4) return false;

  return true;
}

/** Disposes every mesh's geometry and material(s) under `object` — for a
 * rejected `buildFloweryBranch` candidate, or the live one at teardown. */
export function disposeObject3D(object: THREE.Object3D): void {
  object.traverse((child) => {
    if (!(child instanceof THREE.Mesh)) return;
    child.geometry.dispose();
    const mats = Array.isArray(child.material) ? child.material : [child.material];
    mats.forEach((m) => m.dispose());
  });
}
