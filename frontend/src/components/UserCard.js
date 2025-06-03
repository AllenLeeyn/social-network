import Link from "next/link";
import { FaUserCircle } from "react-icons/fa";

export default function UserCard({ user }) {
  return (
    <div className="user-card">
      <Link href={`/profile/${user.uuid}`} className="user-card-link">
        <div className="user-avatar">
          {user.profile_image ? (
            <img
              src={`/frontend-api/image/${user.profile_image}`}
              alt={user.nick_name || "User"}
              className="avatar-img"
            />
          ) : (
            <div className="avatar-placeholder">
              <FaUserCircle size={60} color="#666" />
            </div>
          )}
        </div>
        <div className="user-info">
          <h3 className="user-name">{user.nick_name || "Unknown User"}</h3>
          {user.visibility && (
            <span className={`visibility-badge ${user.visibility}`}>
              {user.visibility}
            </span>
          )}
        </div>
      </Link>
    </div>
  );
}
