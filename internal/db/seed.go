package db

import (
	"context"
	"fmt"
	"github.com/rahuja23/GoBackendProject/internal/store"
	"math/rand/v2"
)

var usernames = [50]string{
	"CyberNinja92",
	"PixelWizard",
	"ShadowRunner",
	"CodeMaster",
	"NeonGamer",
	"ByteHunter",
	"QuantumLeap",
	"FireStorm88",
	"IceBreaker",
	"ThunderBolt",
	"MysticRaven",
	"SteelPhoenix",
	"VoidWalker",
	"StarCrusher",
	"DarkKnight",
	"LightBringer",
	"StormChaser",
	"NightOwl",
	"SilverArrow",
	"GoldenEagle",
	"RedDragon",
	"BlueWave",
	"GreenLantern",
	"PurpleMage",
	"OrangeFlash",
	"YellowSun",
	"PinkPanther",
	"GrayGhost",
	"WhiteWolf",
	"BlackStar",
	"CrimsonBlade",
	"EmeraldShield",
	"SapphireEdge",
	"RubyHeart",
	"DiamondClaw",
	"PearlDiver",
	"CoralReef",
	"OceanTide",
	"MountainPeak",
	"ForestGuard",
	"DesertWind",
	"IceGlacier",
	"VolcanoFire",
	"EarthQuake",
	"SkyLimit",
	"SpaceRanger",
	"TimeKeeper",
	"DreamCatcher",
	"SoulSeeker",
	"MindReader",
}
var postTitles = [20]string{
	"The Rise of AI in Modern Development",
	"10 Essential Tips for Better Code Reviews",
	"Why Remote Work is Here to Stay",
	"Building Scalable Microservices Architecture",
	"The Future of Web Development in 2025",
	"Mastering Docker for Beginners",
	"Understanding Kubernetes Fundamentals",
	"Clean Code Principles Every Developer Should Know",
	"The Art of API Design",
	"Database Optimization Strategies",
	"Cybersecurity Best Practices for Developers",
	"Getting Started with Machine Learning",
	"The Benefits of Test-Driven Development",
	"Cloud Computing vs On-Premise Solutions",
	"Mobile App Development Trends",
	"Version Control with Git: Advanced Techniques",
	"Performance Optimization in Web Applications",
	"The Psychology of User Experience Design",
	"Blockchain Technology Explained Simply",
	"Building Resilient Distributed Systems",
}

var postDescriptions = [20]string{
	"Explore how artificial intelligence is transforming software development practices and what developers need to know to stay ahead in this rapidly evolving landscape.",
	"Discover proven strategies and techniques to conduct more effective code reviews that improve code quality while fostering team collaboration and knowledge sharing.",
	"An in-depth analysis of the remote work revolution, examining its impact on productivity, company culture, and the future of professional collaboration.",
	"Learn how to design and implement microservices architecture that can handle growth, maintain performance, and adapt to changing business requirements.",
	"A comprehensive look at emerging web technologies, frameworks, and trends that will shape the development landscape in the coming year.",
	"Step-by-step guide to containerization with Docker, covering everything from basic concepts to advanced deployment strategies for modern applications.",
	"Master the essentials of Kubernetes orchestration, including pod management, services, deployments, and scaling strategies for containerized applications.",
	"Essential coding principles and practices that lead to maintainable, readable, and efficient code that stands the test of time.",
	"Best practices for designing robust, scalable, and user-friendly APIs that provide excellent developer experience and long-term maintainability.",
	"Advanced techniques for optimizing database performance, including indexing strategies, query optimization, and scaling approaches for high-traffic applications.",
	"Critical security measures every developer should implement to protect applications from common vulnerabilities and emerging threats.",
	"A beginner-friendly introduction to machine learning concepts, tools, and practical applications in software development projects.",
	"Explore the advantages of TDD methodology and learn how to implement it effectively to improve code quality and reduce development time.",
	"Comprehensive comparison of cloud and on-premise solutions, helping you make informed decisions based on cost, security, and scalability factors.",
	"Current trends shaping mobile app development, including cross-platform frameworks, progressive web apps, and emerging user interface patterns.",
	"Advanced Git workflows, branching strategies, and collaboration techniques that improve team productivity and code management.",
	"Proven methods to identify bottlenecks and optimize web application performance for better user experience and reduced server costs.",
	"Understanding user behavior and psychological principles that drive effective UX design decisions and create engaging digital experiences.",
	"Demystifying blockchain technology with practical examples and real-world applications beyond cryptocurrency and digital assets.",
	"Architectural patterns and strategies for building fault-tolerant distributed systems that can handle failures gracefully and maintain high availability.",
}
var postTags = [20]string{
	"AI, Development, Programming, Technology, Innovation",
	"CodeReview, SoftwareDevelopment, BestPractices, TeamWork, Quality",
	"RemoteWork, WorkFromHome, Productivity, Future, Career",
	"Microservices, Architecture, Scalability, Backend, SystemDesign",
	"WebDevelopment, Frontend, JavaScript, React, Trends",
	"Docker, Containerization, DevOps, Deployment, Infrastructure",
	"Kubernetes, Orchestration, Containers, CloudNative, DevOps",
	"CleanCode, Programming, SoftwareEngineering, Quality, Maintainability",
	"API, Design, REST, GraphQL, Backend",
	"Database, Optimization, Performance, SQL, Scaling",
	"Security, Cybersecurity, Development, Privacy, Protection",
	"MachineLearning, AI, DataScience, Python, Algorithms",
	"TDD, Testing, QualityAssurance, Development, Methodology",
	"Cloud, AWS, Azure, OnPremise, Infrastructure",
	"Mobile, iOS, Android, ReactNative, Flutter",
	"Git, VersionControl, Collaboration, Workflow, Development",
	"Performance, Optimization, WebDev, Speed, UserExperience",
	"UX, Design, Psychology, UserInterface, HumanComputer",
	"Blockchain, Cryptocurrency, DeFi, SmartContracts, Web3",
	"DistributedSystems, Architecture, Reliability, Fault-Tolerance, Scaling",
}

var postComments = [20]string{
	"Great insights! I've been using AI coding assistants for months now and the productivity boost is incredible. Have you tried any specific tools you'd recommend?",
	"This is exactly what our team needed. We've been struggling with inconsistent code reviews. The checklist approach sounds promising - definitely implementing this next sprint.",
	"As someone who's been remote for 3 years, I couldn't agree more. The key is setting boundaries and having the right tools. What's your take on hybrid models?",
	"Excellent breakdown of microservices! We're migrating from a monolith right now and this article addresses many of our concerns. The service mesh section was particularly helpful.",
	"Really comprehensive overview of 2025 trends. I'm excited about the WebAssembly developments you mentioned. Do you think it will finally replace JavaScript for performance-critical apps?",
	"Docker changed my development workflow completely. For anyone starting out, I'd add that understanding Docker Compose is equally important. Great tutorial!",
	"Kubernetes seemed overwhelming at first, but this explanation makes it much clearer. The pod lifecycle section was exactly what I needed to understand. Thank you!",
	"These principles should be taught in every CS program. I wish I had learned about clean code earlier in my career. The naming conventions section is gold.",
	"API design is an art form. I learned this the hard way when we had to version our API three times in six months. Consistency is key, as you mentioned.",
	"Database optimization can make or break an application. We saw a 60% performance improvement just by adding proper indexes. The query analysis tips are spot-on.",
	"Security should never be an afterthought. I've seen too many projects rush to production without proper security measures. This checklist is going straight to our team.",
	"Machine learning seemed intimidating before reading this. The practical examples make it much more approachable. Planning to start with the Python libraries you suggested.",
	"TDD transformed how I write code. It feels slower at first, but the confidence it gives you is worth it. The debugging time reduction alone pays for the extra upfront effort.",
	"We're evaluating cloud vs on-premise right now. The cost breakdown in your article is really helpful. Have you seen any changes in pricing models recently?",
	"Mobile development is evolving so fast. Flutter has been a game-changer for our cross-platform needs. What's your experience with React Native vs Flutter?",
	"Git workflows can make or break team productivity. We switched to the GitFlow model you described and our merge conflicts dropped by 80%. Highly recommend it.",
	"Performance optimization is crucial for user retention. We followed similar strategies and reduced our bounce rate significantly. The lazy loading technique works wonders.",
	"UX psychology is fascinating. Understanding user behavior patterns has helped us design much more intuitive interfaces. The cognitive load section was particularly insightful.",
	"Blockchain is finally moving beyond just cryptocurrency. The smart contract examples you provided show real practical applications. Excited to experiment with Web3 development.",
	"Building resilient systems is challenging but essential. We implemented circuit breakers after reading about them here and it saved us during a recent traffic spike.",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(50)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			fmt.Printf("The error occured while seeding user %s: %v\n", user.Username, err)
		}
	}
	posts := generatePost(20, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			fmt.Printf("The error occured while seeding post %s: %v\n", post.Title, err)
		}
	}
	comments := generateComments(20, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			fmt.Printf("The error occured while seeding comments %s: %v\n", comment.Content, err)
		}
	}
	fmt.Printf("Seeding Completed Successfully!\n")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)
	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[rand.IntN(len(users))] + fmt.Sprintf("%d", i),
			Email:    usernames[rand.IntN(len(users))] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123123",
		}

	}

	return users
}
func generatePost(count int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, count)
	for i := 0; i < count; i++ {
		user := users[rand.IntN(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   postTitles[rand.IntN(len(postTitles))],
			Content: postDescriptions[rand.IntN(len(postDescriptions))],
			Tags: []string{
				postTags[rand.IntN(len(postTitles))],
				postTags[rand.IntN(len(postTitles))],
			},
		}
	}
	return posts
}

func generateComments(count int, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, count)
	for i := 0; i < count; i++ {
		comments[i] = &store.Comment{
			PostID:  posts[i].ID,
			UserID:  posts[i].UserID,
			Content: postComments[rand.IntN(len(postComments))],
		}
	}
	return comments
}
