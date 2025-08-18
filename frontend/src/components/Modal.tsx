import { motion, AnimatePresence } from 'framer-motion'
interface ModalProps {
  isOpen: boolean
  onClose: () => void
  children: React.ReactNode
}

export const Modal: React.FC<ModalProps> = ({ isOpen, onClose, children }) => {
  const modalVariants = {
    hidden: { y: '100%' },
    visible: {
      y: '25%', // Slide up to 75% of the screen height
      transition: { duration: 0.3, ease: 'easeInOut' },
    },
    exit: { y: '100%', transition: { duration: 0.2 } },
  }

  const backdropVariants = {
    hidden: { opacity: 0 },
    visible: { opacity: 0.5, transition: { duration: 0.3 } },
    exit: { opacity: 0, transition: { duration: 0.2 } },
  }

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          <motion.div
            variants={backdropVariants}
            initial="hidden"
            animate="visible"
            exit="exit"
            className="fixed inset-0 bg-black z-50"
            onClick={onClose} // Close modal when backdrop is clicked.
          />
          <motion.div
            variants={modalVariants}
            initial="hidden"
            animate="visible"
            exit="exit"
            className="fixed bottom-0 left-0 right-0 bg-white rounded-t-[40px] z-50 overflow-auto min-h-[75vh]" // Ensure it doesn't exceed 75% of viewport height
            style={{
              WebkitOverflowScrolling: 'touch', // Enable smooth scrolling on iOS
            }}
          >
            {children}
          </motion.div>
        </>
      )}
    </AnimatePresence>
  )
}
