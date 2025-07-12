import { useEffect, useState } from "react";
import { SendData } from "../../../../../../utils/sendData.js";
import Image from "next/image.js";

export default function Members({ groupId }) {
    const [members, setMembers] = useState([]);
    const [hasError, setHasError] = useState(false);

    useEffect(() => {
        async function fetchMembers() {
            if (!groupId) return;

            try {
                const res = await SendData(`/api/v1/get/groupMembers`, groupId);

                if (res.ok) {
                    const memberData = await res.json();

                    // Access the `members` array inside the response object
                    if (memberData.members && Array.isArray(memberData.members)) {
                        setMembers(memberData.members);  // Set the actual array of members
                    } else {
                        setHasError(true);
                        console.error("No 'members' array found in response:", memberData);
                    }
                } else {
                    setHasError(true);
                    console.error("Failed to fetch members: ", res.status);
                }
            } catch (err) {
                setHasError(true);
                console.error("Error fetching members:", err);
            }
        }

        fetchMembers();
    }, [groupId]);

    if (hasError) {
        return <div>Error loading members.</div>;
    }

    return (
        <div>
            {members && members.length > 0 ? (
                members.map((member, index) => (
                    <div key={index}>
                        <Image src={member?.Avatar?.Valid ? member?.Avatar?.String : "/iconMale.png"}
                            alt={member.Name}
                            width={50}
                            height={50}
                        />
                        <p>{member.Name}</p>
                    </div>
                ))
            ) : (
                <div>No members available.</div>
            )}
        </div>
    );
}
