import React from 'react';
import { FaChartBar, FaEdit, FaUsers, FaUpload } from 'react-icons/fa';

function CreatorDashboard() {
  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Creator Dashboard</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <DashboardCard icon={<FaChartBar />} title="Analytics" value="10.5k" subtitle="Total Views" />
        <DashboardCard icon={<FaEdit />} title="Content" value="25" subtitle="Published Items" />
        <DashboardCard icon={<FaUsers />} title="Followers" value="2.7k" subtitle="Total Followers" />
        <DashboardCard icon={<FaUpload />} title="Uploads" value="15" subtitle="Pending Uploads" />
      </div>

      <div className="bg-white shadow-md rounded-lg p-6 mb-8">
        <h2 className="text-xl font-semibold mb-4">Recent Activity</h2>
        <ul className="space-y-4">
          <ActivityItem 
            title="New follower" 
            description="John Doe started following you" 
            time="2 hours ago" 
          />
          <ActivityItem 
            title="Content published" 
            description="Your new item 'Summer Collection' is now live" 
            time="1 day ago" 
          />
          <ActivityItem 
            title="Comment received" 
            description="Alice left a comment on your 'Spring Trends' post" 
            time="3 days ago" 
          />
        </ul>
      </div>

      <div className="bg-white shadow-md rounded-lg p-6">
        <h2 className="text-xl font-semibold mb-4">Quick Actions</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <ActionButton icon={<FaUpload />} text="Upload New Content" />
          <ActionButton icon={<FaEdit />} text="Edit Profile" />
          <ActionButton icon={<FaUsers />} text="Engage with Followers" />
        </div>
      </div>
    </div>
  );
}

function DashboardCard({ icon, title, value, subtitle }) {
  return (
    <div className="bg-white shadow-md rounded-lg p-6">
      <div className="flex items-center justify-between mb-4">
        <div className="text-3xl text-blue-500">{icon}</div>
        <h3 className="text-lg font-semibold">{title}</h3>
      </div>
      <p className="text-2xl font-bold mb-1">{value}</p>
      <p className="text-sm text-gray-600">{subtitle}</p>
    </div>
  );
}

function ActivityItem({ title, description, time }) {
  return (
    <li className="flex items-center">
      <div className="bg-blue-100 rounded-full p-2 mr-4">
        <FaUsers className="text-blue-500" />
      </div>
      <div>
        <h4 className="font-semibold">{title}</h4>
        <p className="text-sm text-gray-600">{description}</p>
        <p className="text-xs text-gray-400">{time}</p>
      </div>
    </li>
  );
}

function ActionButton({ icon, text }) {
  return (
    <button className="flex items-center justify-center bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition-colors">
      <span className="mr-2">{icon}</span>
      {text}
    </button>
  );
}

export default CreatorDashboard;