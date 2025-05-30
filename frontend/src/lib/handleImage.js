const API_URL = "/frontend-api";

export const handleImage = async (files) => {
  const MAX_SIZE_MB = 5; // e.g. 5 MB
  const MAX_SIZE_BYTES = MAX_SIZE_MB * 1024 * 1024;

  if (!files || files.length === 0) throw new Error("No file selected.");

  const formData = new FormData();
  files.forEach((file) => {
    if (file.size > MAX_SIZE_BYTES) {
        throw new Error(`File size exceeds ${MAX_SIZE_MB}MB limit.`);
    }
    formData.append("fileAttachments", file);
  });

  const response = await fetch(`${API_URL}/imageUpload`, {
      method: "POST",
      body: formData,
  });

  if (!response.ok) {
      throw new Error("Failed to upload image");
  }

  const data = await response.json();
  return data.data;
};
