import styles from "../profile.module.css";

export default function CreateEvent({ groupId }) {
  const handleSubmit = async (e) => {
    e.preventDefault();
    const form = e.target;

    const title = form.title.value.trim();
    const description = form.description.value.trim();
    const datetime = form.datetime.value;

    const options = Array.from(form.options.selectedOptions).map(opt => opt.value);

    if (options.length < 2) {
      alert("Please select at least 2 options.");
      return;
    }

    const payload = {
      groupId,
      title,
      description,
      datetime,
      options
    };

    console.log(payload);
    // await SendData("/api/v1/events/create", payload);
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit} style={{ width: "97%",marginTop:"10px",marginLeft:"-10px" }}>
      <input type="text" name="title" placeholder="Event Title" required />
      <textarea name="description" placeholder="Event Description" required />
      <input type="datetime-local" name="datetime" required />

      <h4>Options (Select at least two)</h4>
      <select name="options" required>
        <option value="Going">Going</option>
        <option value="Not Going">Not Going</option>
      </select>

      <button type="submit">Create Event</button>
    </form>
  );
}
