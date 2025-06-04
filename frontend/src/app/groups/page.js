'use client';

import React, { useState } from 'react';
import SidebarSection from '../../components/SidebarSection';
import UsersList from '../../components/UsersList';
import Modal from '../../components/Modal';
import GroupFilterList from '../../components/groups/GroupFilterList';
import GroupList from '../../components/groups/GroupList';
import GroupDetail from '../../components/groups/GroupDetail';
import GroupInvitationList from '../../components/groups/GroupInvitationList';
import CreateGroupForm from '../../components/groups/CreateGroupForm';

import '../../styles/groups/FilterList.css'
import './groups.css'

const groupFilters = [
  { key: 'my_groups', label: 'My Groups' },
  { key: 'discover', label: 'Discover' },
  { key: 'invitations', label: 'Invitations' }
];

export default function GroupsPage() {
  const [selectedFilter, setSelectedFilter] = useState('my_groups');
  const [selectedGroup, setSelectedGroup] = useState(null);
  const [isCreateModalOpen, setCreateModalOpen] = useState(false);
  const [groupsRefreshKey, setGroupsRefreshKey] = useState(0);


  const handleSelectGroup = (group) => {
  if (selectedGroup && selectedGroup.id === group.id) {
    setSelectedGroup(null); // Deselect if same group is clicked
  } else {
    setSelectedGroup(group); // Select new group
  }
};

  return (
    <main>
      <div className='groups-page-layout'>

        {/* Left Sidebar */}
        <aside className='sidebar left-sidebar'>
          <SidebarSection title='Groups'>
            <GroupFilterList
              filters={groupFilters}
              selectedFilter={selectedFilter}
              onSelect={setSelectedFilter}
            />
          </SidebarSection>
        </aside>

        {/* Main Content */}
        <section className='main-feed group-section'>

        {/* Create Group Button */}
        <div style={{ margin: '1rem 0', textAlign: 'right' }}>
          <button onClick={() => setCreateModalOpen(true)}>
            + Create Group
          </button>
        </div>
          {selectedFilter === 'my_groups' && (
            <GroupList key={groupsRefreshKey} onSelectGroup={handleSelectGroup} />
          )}
          {selectedFilter === 'discover' && (
            <GroupList discover onSelectGroup={handleSelectGroup} />
          )}
          {selectedFilter === 'invitations' && (
            <GroupInvitationList />
          )}
          {selectedGroup && (
            <GroupDetail group={selectedGroup} />
          )}

          {isCreateModalOpen && (
            <Modal title="Create Group" onClose={() => setCreateModalOpen(false)}>
              <CreateGroupForm 
                onSuccess={() => {
                  setCreateModalOpen(false)
                  setGroupsRefreshKey(k => k+1);
                }} 
              />
            </Modal>
          )}
        </section>

        {/* Right Sidebar */}
        <aside className="sidebar right-sidebar">
          <SidebarSection title="All Users">
            <UsersList />
          </SidebarSection>
        </aside>
      </div>
    </main>
  );
}
