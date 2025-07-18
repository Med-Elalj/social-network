import styles from "../profile.module.css";
import { SendData } from "@/app/sendData.js";
import { useNotification } from "@/app/context/NotificationContext";

export default function CreateEvent({ groupId, setActiveSection }) {
  const { showNotification } = useNotification();

  const handleSubmit = async (e) => {
    e.preventDefault();
    const form = e.target;

    const title = form.title.value.trim();
    const description = form.description.value.trim();
    const datetime = form.datetime.value;

    const selectedTime = new Date(datetime).getTime();
    const oneYearFromNow = Date.now() + 1 * 365 * 24 * 60 * 60 * 1000;

    if (selectedTime <= Date.now()) {
      showNotification("Date and time must be in the future", "error");
      return;
    }

    if (selectedTime > oneYearFromNow) {
      showNotification("Date and time cannot be more than 1 year in the future", "error");
      return;
    }
    
    if (title === "" || description === "" || datetime === "") {
      showNotification("All fields are required", "error");
      return;
    }

    const payload = {
      group_id: groupId,
      title,
      description,
      time: datetime,
    };

    const response = await SendData("/api/v1/set/eventCreation", payload);

    if (response.status !== 200) {
      const errorBody = await response.json();
      console.log(errorBody);
    } else {
      showNotification("Event created successfully!", "success");
      setActiveSection("posts");
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit} style={{ width: "97%", marginTop: "10px", marginLeft: "-10px" }}>
      <label htmlFor="title">Title</label>
      <input type="text" name="title" placeholder="Event Title" required />

      <label htmlFor="description">Description</label>
      <textarea name="description" placeholder="Event Description" required />

      <label htmlFor="datetime">Date and Time</label>
      <input type="datetime-local" name="datetime" required />

      <button type="submit">Create Event</button>
    </form>
  );
}
