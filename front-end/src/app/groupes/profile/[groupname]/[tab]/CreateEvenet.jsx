import styles from "../profile.module.css";

export default function CreateEvent({ groupId }) {
  const handleSubmit = async (e) => {
    e.preventDefault();
    const form = e.target;

    const title = form.title.value.trim();
    const description = form.description.value.trim();
    const datetime = form.datetime.value;

    const payload = {
      groupId,
      title,
      description,
      datetime,
    };

    console.log(payload);
    // await SendData("/api/v1/events/create", payload);
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
