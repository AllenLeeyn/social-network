const API_URL = "/frontend-api";

export const handleImage = async (file) => {
  const MAX_SIZE_MB = 5; // e.g. 5 MB
  const MAX_SIZE_BYTES = MAX_SIZE_MB * 1024 * 1024;


  if (!file) throw new Error("No file selected.");
  if (file.size > MAX_SIZE_BYTES) {
    throw new Error(`File size exceeds ${MAX_SIZE_MB}MB limit.`);
  }

  const formData = new FormData();
  formData.append("fileAttachments", file);

  const response = await fetch(`${API_URL}/imageUpload`, {
      method: "POST",
      body: formData,
  });

  if (!response.ok) {
      throw new Error("Failed to upload image");
  }

  const data = await response.json();
  return Object.values(data.data)[0];
};
