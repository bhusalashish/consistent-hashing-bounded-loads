# Quick Copy-Paste Prompt for AI Presentation Tools

Copy this prompt into Gamma, Beautiful.ai, ChatGPT, or any AI presentation generator:

---

**Create a professional 20-slide academic presentation titled "Consistent Hashing with Bounded Loads: An Interactive Visualization Tool" for a Distributed Systems course (CMPE 273).**

**Slide Structure:**

1. **Title Slide**: "Consistent Hashing with Bounded Loads: An Interactive Visualization Tool" - CMPE 273 Distributed Systems

2. **Problem Statement**: Load balancing challenges in distributed systems - traditional hashing fails on node changes, need for consistent mapping, load imbalance causes hotspots

3. **Consistent Hashing Overview**: Hash ring concept - nodes on circle, keys hash to positions, assign to first node clockwise, minimal key movement

4. **Load Imbalance Problem**: Uneven distribution (some nodes get 2-3x more keys), causes hotspots, example bar chart showing imbalance

5. **Four Algorithms**: 
   - Ring Consistent Hash (RingCH) - basic with virtual nodes
   - Jump Consistent Hashing - O(log n), no ring
   - Maglev - Google's production load balancer
   - CH-BL - Consistent Hashing with Bounded Loads (our focus)

6. **CH-BL Deep Dive**: 
   - Capacity formula: C = ⌈c × (ExpectedKeys / numNodes)⌉
   - Load factor c = 1.25 (25% over-provisioning)
   - Walking algorithm: hash → walk clockwise → find capacity
   - Two-choice fallback for long walks

7. **System Architecture**: Go backend (HTTP/JSON APIs) + React/TypeScript/D3.js frontend, real-time visualization

8. **Key Features**: Visual hash ring, real-time animations, statistics dashboard, algorithm comparison, educational tooltips, CH-BL configuration

9. **Demo Setup**: "Live Demo - Let's see it in action!"

10. **Demo Part 1**: Adding/removing nodes, showing key movements and churn statistics

11. **Demo Part 2**: Algorithm comparison - switching between algorithms, comparing churn percentages

12. **Demo Part 3**: CH-BL capacity enforcement - showing capacity limits, key walking, adjusting load factor

13. **Experimental Results**: 
    - Churn comparison: RingCH ~25%, Jump ~20%, Maglev low, CH-BL similar to RingCH
    - Load distribution: RingCH can have 2-3x imbalance, CH-BL guaranteed within bounds

14. **Real-World Applications**: CDNs, distributed databases, caching systems (Memcached, Redis), load balancers (Google Maglev, AWS ELB)

15. **Technical Challenges**: CH-BL state management, capacity calculation, D3.js animations, real-time statistics tracking

16. **Key Takeaways**: Consistent hashing minimizes movement, load imbalance is real problem, CH-BL provides guaranteed bounds, visualization helps understanding

17. **Future Work**: Weighted nodes, multi-dimensional metrics, performance benchmarking, real system integration

18. **Q&A**: Thank you slide with contact info and GitHub link

**Design Requirements:**
- Google Material Design colors (primary blue #1a73e8, green #34a853, yellow #fbbc04, red #ea4335)
- Clean white background with subtle gradients
- Google Sans/Roboto fonts
- Visual diagrams: hash rings, nodes, keys, comparison charts
- Smooth transitions and animations
- Professional academic style

**Each slide should include relevant visuals: diagrams, charts, code snippets, or screenshots. The presentation should be engaging and suitable for a 20-minute academic presentation with a 5-minute live demo component.**

---

## Alternative: One-Paragraph Version

**Create a 20-slide academic presentation on "Consistent Hashing with Bounded Loads: An Interactive Visualization Tool" covering: problem statement of load balancing in distributed systems, consistent hashing overview with hash ring visualization, comparison of four algorithms (RingCH, Jump, Maglev, CH-BL), deep dive into CH-BL capacity bounds and walking algorithm, system architecture (Go backend + React/D3.js frontend), interactive features showcase, 5-minute live demo walkthrough showing node operations and algorithm comparison, experimental results comparing churn and load distribution, real-world applications (CDNs, databases, caching), technical challenges, key takeaways, and future work. Use Google Material Design color scheme (#1a73e8 primary blue), clean professional layout, visual diagrams of hash rings and nodes, comparison charts, smooth animations, and Google Sans/Roboto typography. Suitable for a 20-minute Distributed Systems course presentation.**

