const express = require('express');
const bcrypt = require('bcrypt');
const jwt = require('jsonwebtoken');
const QRCode = require('qrcode');
const cors = require('cors');
const rateLimit = require('express-rate-limit');
const helmet = require('helmet');
const { body, validationResult, param } = require('express-validator');

// Initialize Express app
const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(helmet());
app.use(cors());
app.use(express.json({ limit: '10mb' }));

// Rate limiting
const limiter = rateLimit({
    windowMs: 15 * 60 * 1000, // 15 minutes
    max: 100 // limit each IP to 100 requests per windowMs
});
app.use(limiter);

// Mock database (replace with actual database in production)
const users = new Map();
const nfts = new Map();
const organizations = new Map();

// JWT Secret (use environment variable in production)
const JWT_SECRET = process.env.JWT_SECRET || 'your-secret-key';

// Hyperledger Fabric SDK integration (mock for this example)
class HyperledgerService {
    static async queryNFT(nftId) {
        // Mock implementation - replace with actual Hyperledger Fabric SDK calls
        return {
            id: nftId,
            title: 'Bachelor of Computer Science',
            studentId: 'student123',
            universityId: 'university456',
            issueDate: new Date().toISOString(),
            metadata: {
                gpa: '3.8',
                honors: 'Magna Cum Laude'
            }
        };
    }

    static async verifyNFT(nftId) {
        // Mock implementation - replace with actual verification logic
        return {
            isValid: true,
            timestamp: new Date().toISOString(),
            blockHash: 'abc123...',
            transactionId: 'tx456...'
        };
    }

    static async getAllNFTsByOrganization(orgId) {
        // Mock implementation
        return [
            {
                id: 'nft1',
                title: 'Bachelor of Computer Science',
                studentId: 'student123',
                issueDate: new Date().toISOString()
            }
        ];
    }
}

// Middleware for authentication
const authenticateToken = (req, res, next) => {
    const authHeader = req.headers['authorization'];
    const token = authHeader && authHeader.split(' ')[1];

    if (!token) {
        return res.status(401).json({ error: 'Access token required' });
    }

    jwt.verify(token, JWT_SECRET, (err, user) => {
        if (err) {
            return res.status(403).json({ error: 'Invalid token' });
        }
        req.user = user;
        next();
    });
};

// Middleware for organization authentication
const authenticateOrganization = (req, res, next) => {
    if (req.user.role !== 'organization') {
        return res.status(403).json({ error: 'Organization access required' });
    }
    next();
};

// Validation middleware
const validateRequest = (req, res, next) => {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
        return res.status(400).json({ errors: errors.array() });
    }
    next();
};

// AUTHENTICATION ROUTES

// Student Registration
app.post('/api/auth/register/student', [
    body('email').isEmail().normalizeEmail(),
    body('password').isLength({ min: 8 }),
    body('studentId').notEmpty(),
    body('firstName').notEmpty(),
    body('lastName').notEmpty(),
    body('university').notEmpty()
], validateRequest, async (req, res) => {
    try {
        const { email, password, studentId, firstName, lastName, university } = req.body;

        // Check if user already exists
        if (users.has(email)) {
            return res.status(400).json({ error: 'User already exists' });
        }

        // Hash password
        const hashedPassword = await bcrypt.hash(password, 10);

        // Create user
        const user = {
            email,
            password: hashedPassword,
            studentId,
            firstName,
            lastName,
            university,
            role: 'student',
            createdAt: new Date().toISOString(),
            verified: false
        };

        users.set(email, user);

        // Generate JWT
        const token = jwt.sign(
            { email, role: 'student', studentId },
            JWT_SECRET,
            { expiresIn: '24h' }
        );

        res.status(201).json({
            message: 'Student registered successfully',
            token,
            user: {
                email: user.email,
                studentId: user.studentId,
                firstName: user.firstName,
                lastName: user.lastName,
                university: user.university,
                role: user.role
            }
        });
    } catch (error) {
        res.status(500).json({ error: 'Internal server error' });
    }
});

// Organization Registration
app.post('/api/auth/register/organization', [
    body('email').isEmail().normalizeEmail(),
    body('password').isLength({ min: 8 }),
    body('organizationName').notEmpty(),
    body('organizationId').notEmpty(),
    body('contactPerson').notEmpty()
], validateRequest, async (req, res) => {
    try {
        const { email, password, organizationName, organizationId, contactPerson } = req.body;

        if (organizations.has(email)) {
            return res.status(400).json({ error: 'Organization already exists' });
        }

        const hashedPassword = await bcrypt.hash(password, 10);

        const organization = {
            email,
            password: hashedPassword,
            organizationName,
            organizationId,
            contactPerson,
            role: 'organization',
            createdAt: new Date().toISOString(),
            verified: false
        };

        organizations.set(email, organization);

        const token = jwt.sign(
            { email, role: 'organization', organizationId },
            JWT_SECRET,
            { expiresIn: '24h' }
        );

        res.status(201).json({
            message: 'Organization registered successfully',
            token,
            organization: {
                email: organization.email,
                organizationName: organization.organizationName,
                organizationId: organization.organizationId,
                contactPerson: organization.contactPerson,
                role: organization.role
            }
        });
    } catch (error) {
        res.status(500).json({ error: 'Internal server error' });
    }
});

// Login
app.post('/api/auth/login', [
    body('email').isEmail().normalizeEmail(),
    body('password').notEmpty()
], validateRequest, async (req, res) => {
    try {
        const { email, password } = req.body;

        // Check in both users and organizations
        const user = users.get(email) || organizations.get(email);

        if (!user) {
            return res.status(401).json({ error: 'Invalid credentials' });
        }

        // Verify password
        const isValidPassword = await bcrypt.compare(password, user.password);
        if (!isValidPassword) {
            return res.status(401).json({ error: 'Invalid credentials' });
        }

        // Generate JWT
        const tokenPayload = {
            email: user.email,
            role: user.role,
            ...(user.role === 'student' ? { studentId: user.studentId } : { organizationId: user.organizationId })
        };

        const token = jwt.sign(tokenPayload, JWT_SECRET, { expiresIn: '24h' });

        res.json({
            message: 'Login successful',
            token,
            user: {
                email: user.email,
                role: user.role,
                ...(user.role === 'student' ? {
                    studentId: user.studentId,
                    firstName: user.firstName,
                    lastName: user.lastName,
                    university: user.university
                } : {
                    organizationName: user.organizationName,
                    organizationId: user.organizationId,
                    contactPerson: user.contactPerson
                })
            }
        });
    } catch (error) {
        res.status(500).json({ error: 'Internal server error' });
    }
});

// STUDENT ROUTES

// Get student's NFTs
app.get('/api/student/nfts', authenticateToken, async (req, res) => {
    try {
        if (req.user.role !== 'student') {
            return res.status(403).json({ error: 'Student access required' });
        }

        // Mock implementation - replace with actual Hyperledger query
        const studentNFTs = [
            {
                id: 'nft1',
                title: 'Bachelor of Computer Science',
                issuer: 'University of Technology',
                issueDate: '2024-05-15',
                metadata: {
                    gpa: '3.8',
                    honors: 'Magna Cum Laude'
                },
                verified: true
            },
            {
                id: 'nft2',
                title: 'Certificate in Data Science',
                issuer: 'University of Technology',
                issueDate: '2024-03-10',
                metadata: {
                    grade: 'A'
                },
                verified: true
            }
        ];

        res.json({
            nfts: studentNFTs,
            count: studentNFTs.length
        });
    } catch (error) {
        res.status(500).json({ error: 'Internal server error' });
    }
});

// Get specific NFT details
app.get('/api/student/nfts/:nftId', [
    param('nftId').notEmpty()
], validateRequest, authenticateToken, async (req, res) => {
    try {
        if (req.user.role !== 'student') {
            return res.status(403).json({ error: 'Student access required' });
        }

        const nftId = req.params.nftId;

        // Query from Hyperledger
        const nft = await HyperledgerService.queryNFT(nftId);

        if (!nft) {
            return res.status(404).json({ error: 'NFT not found' });
        }

        res.json({ nft });
    } catch (error) {
        res.status(500).json({ error: 'Internal server error' });
    }
});

// ORGANIZATION ROUTES

// Get all NFTs issued by organization
app.get('/api/organization/nfts', authenticateToken, authenticateOrganization, async (req, res) => {
    try {
        const organizationId = req.user.organizationId;

        // Query all NFTs issued by this organization
        const nfts = await HyperledgerService.getAllNFTsByOrganization(organizationId);

        res.json({
            nfts,
            count: nfts.length,
            organizationId
        });
    } catch (error) {
        res.status(500).json({ error: 'Internal server error' });
    }
});

// Get NFT traceability
app.get('/api/organization/nfts/:nftId/trace', [
    param('nftId').notEmpty()
], validateRequest, authenticateToken, authenticateOrganization, async (req, res) => {
    try {
        const nftId = req.params.nftId;

        // Mock traceability data - replace with actual blockchain query
        const traceData = {
            nftId,
            creationDate: '2024-05-15T10:30:00Z',
            creator: req.user.organizationId,
            transactions: [
                {
                    txId: 'tx123',
                    timestamp: '2024-05-15T10:30:00Z',
                    action: 'CREATE',
                    blockHeight: 1001
                },
                {
                    txId: 'tx124',
                    timestamp: '2024-05-15T10:31:00Z',
                    action: 'ASSIGN',
                    recipient: 'student123',
                    blockHeight: 1002
                }
            ],
            currentOwner: 'student123',
            verified: true
        };

        res.json({ traceability: traceData });
    } catch (error) {
        res.status(500).json({ error: 'Internal server error' });
    }
});

// PUBLIC VERIFICATION ROUTES

// Verify NFT (for external agents)
app.get('/api/verify/:nftId', [
    param('nftId').notEmpty()
], validateRequest, async (req, res) => {
    try {
        const nftId = req.params.nftId;

        // Verify NFT on blockchain
        const verification = await HyperledgerService.verifyNFT(nftId);
        const nftData = await HyperledgerService.queryNFT(nftId);

        res.json({
            nftId,
            isValid: verification.isValid,
            nft: nftData,
            verification: {
                timestamp: verification.timestamp,
                blockHash: verification.blockHash,
                transactionId: verification.transactionId
            }
        });
    } catch (error) {
        res.status(500).json({ error: 'Internal server error' });
    }
});

// Generate QR code for NFT
app.get('/api/qr/:nftId', [
    param('nftId').notEmpty()
], validateRequest, async (req, res) => {
    try {
        const nftId = req.params.nftId;

        // Create verification URL
        const verificationUrl = `${req.protocol}://${req.get('host')}/api/verify/${nftId}`;

        // Generate QR code
        const qrCode = await QRCode.toDataURL(verificationUrl);

        res.json({
            nftId,
            qrCode,
            verificationUrl
        });
    } catch (error) {
        res.status(500).json({ error: 'Internal server error' });
    }
});

// Get QR code as image
app.get('/api/qr/:nftId/image', [
    param('nftId').notEmpty()
], validateRequest, async (req, res) => {
    try {
        const nftId = req.params.nftId;

        const verificationUrl = `${req.protocol}://${req.get('host')}/api/verify/${nftId}`;

        // Generate QR code as PNG buffer
        const qrBuffer = await QRCode.toBuffer(verificationUrl, {
            type: 'png',
            width: 300,
            margin: 2
        });

        res.set('Content-Type', 'image/png');
        res.send(qrBuffer);
    } catch (error) {
        res.status(500).json({ error: 'Internal server error' });
    }
});

// UTILITY ROUTES

// Health check
app.get('/api/health', (req, res) => {
    res.json({
        status: 'OK',
        timestamp: new Date().toISOString(),
        version: '1.0.0'
    });
});

// Get user profile
app.get('/api/profile', authenticateToken, (req, res) => {
    const userEmail = req.user.email;
    const user = users.get(userEmail) || organizations.get(userEmail);

    if (!user) {
        return res.status(404).json({ error: 'User not found' });
    }

    const { password, ...userProfile } = user;
    res.json({ profile: userProfile });
});

// Error handling middleware
app.use((err, req, res, next) => {
    console.error(err.stack);
    res.status(500).json({ error: 'Something went wrong!' });
});

// 404 handler
app.use((req, res) => {
    res.status(404).json({ error: 'Route not found' });
});

// Start server
app.listen(PORT, () => {
    console.log(`NFT Wallet API server running on port ${PORT}`);
    console.log(`Health check: http://localhost:${PORT}/api/health`);
});

module.exports = app;