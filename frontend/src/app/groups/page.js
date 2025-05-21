import styles from './groups.css';

export default function GroupsPage() {
  const groups = mockGroups;

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Groups</h1>
      <div className={styles.groupList}>
        {groups.map(group => (
          <div key={group.id} className={styles.groupCard}>
            <h2 className={styles.groupName}>{group.name}</h2>
            <p className={styles.groupDescription}>{group.description}</p>
            <p className={styles.groupInfo}>Members: {group.memberCount}</p>
            <span className={styles.groupTag}>
              {group.isPublic ? "Public" : "Private"}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}


const mockGroups = [
  {
    id: 1,
    name: "Tech Enthusiasts",
    description: "A group for sharing tech tips and news.",
    memberCount: 142,
    isPublic: true
  },
  {
    id: 2,
    name: "Book Club",
    description: "Discuss your favorite books every week.",
    memberCount: 87,
    isPublic: true
  },
  {
    id: 3,
    name: "Fitness Group",
    description: "Workout routines and motivation.",
    memberCount: 65,
    isPublic: true
  },
  {
    id: 4,
    name: "Cooking Masters",
    description: "Recipes and cooking techniques.",
    memberCount: 53,
    isPublic: true
  }
];
