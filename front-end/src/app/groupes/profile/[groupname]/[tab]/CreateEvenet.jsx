import styles from "../profile.module.css";
import { SendData } from "@/app/sendData.js";

export default function CreateEvent({ groupId }) {
  const handleSubmit = async (e) => {
    e.preventDefault();
    const form = e.target;

    const title = form.title.value.trim();
    const description = form.description.value.trim();
    const datetime = form.datetime.value;

    const payload = {
      group_id: groupId,
      title: title,
      description: description,
      time: datetime,
    };

    console.log("front event data : ", payload);
    const response = await SendData("/api/v1/set/eventCreation", payload);
    const Body = await response.json();

    if (response.status !== 200) {
      const errorBody = await response.json();
      console.log(errorBody);
    } else {
      console.log("✅ Event created!");
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
